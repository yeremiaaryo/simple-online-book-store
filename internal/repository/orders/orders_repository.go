package orders

import (
	"context"
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
