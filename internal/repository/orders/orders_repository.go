package orders

import (
	"context"
	"github.com/lib/pq"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/orders"
	"github.com/yeremiaaryo/gotu-assignment/pkg/internalsql"
	"time"
)

type repository struct {
	masterDB internalsql.MasterDB
	slaveDB  internalsql.SlaveDB
}

func New(masterDB internalsql.MasterDB, slaveDB internalsql.SlaveDB) *repository {
	r := repository{
		masterDB: masterDB,
		slaveDB:  slaveDB,
	}

	return &r
}

func (r *repository) InsertOrder(ctx context.Context, order orders.CreateOrderRequest) (*orders.CreateOrderResponse, error) {
	tx, err := r.masterDB.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	createdAt := time.Now().UnixMilli()
	updatedAt := createdAt

	stmtOrder, err := tx.PreparexContext(ctx, tx.Rebind(insertOrderQuery))
	if err != nil {
		return nil, err
	}
	defer stmtOrder.Close()

	var orderID int64
	err = stmtOrder.QueryRowContext(ctx, order.UserID, order.TotalAmount, orders.OrderStatusNew, createdAt, updatedAt).Scan(&orderID)
	if err != nil {
		return nil, err
	}

	stmtOrderItem, err := tx.PreparexContext(ctx, tx.Rebind(insertOrderItemQuery))
	if err != nil {
		return nil, err
	}
	defer stmtOrderItem.Close()

	for _, item := range order.Items {
		_, err = stmtOrderItem.ExecContext(ctx, orderID, item.BookID, item.Quantity, item.Price, createdAt, updatedAt)
		if err != nil {
			return nil, err
		}
	}

	response := &orders.CreateOrderResponse{
		OrderID: orderID,
		Status:  orders.OrderStatusNew.String(),
	}

	return response, tx.Commit()
}

func (r *repository) GetOrdersByUserID(ctx context.Context, userID int64, limit, offset int) ([]orders.History, error) {
	rebindQuery := r.slaveDB.Rebind(getOrderHistoryByUserID)
	stmt, err := r.slaveDB.PreparexContext(ctx, rebindQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryxContext(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordersList []orders.History
	for rows.Next() {
		var order orders.History
		err = rows.Scan(&order.ID, &order.TotalAmount, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		ordersList = append(ordersList, order)
	}

	if len(ordersList) == 0 {
		return nil, nil
	}

	orderIDs := make([]int64, len(ordersList))
	for i, order := range ordersList {
		orderIDs[i] = order.ID
	}

	rebindItemQuery := r.slaveDB.Rebind(getItemsQuery)
	stmtItem, err := r.slaveDB.PreparexContext(ctx, rebindItemQuery)
	if err != nil {
		return nil, err
	}
	defer stmtItem.Close()

	itemsRows, err := stmtItem.QueryxContext(ctx, pq.Array(orderIDs))
	if err != nil {
		return nil, err
	}
	defer itemsRows.Close()

	// Map items to their respective orders
	itemsMap := make(map[int64][]orders.ItemHistory)
	for itemsRows.Next() {
		var (
			item    orders.ItemHistory
			orderID int64
		)

		err = itemsRows.Scan(&item.ID, &orderID, &item.BookID, &item.Quantity, &item.Price)
		if err != nil {
			return nil, err
		}
		itemsMap[orderID] = append(itemsMap[orderID], item)
	}

	// Attach items to their orders
	for i, order := range ordersList {
		ordersList[i].Items = itemsMap[order.ID]
	}

	return ordersList, nil
}
