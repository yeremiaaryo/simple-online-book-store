package orders

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/orders"
	"github.com/yeremiaaryo/gotu-assignment/pkg/internalsql"
	"reflect"
	"testing"
)

func Test_repository_InsertOrder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	masterDB := internalsql.NewMasterDB(db, "sqlmock")
	defer func() {
		_ = db.Close()
	}()

	insertOrderQueryTest := masterDB.Rebind(`
        INSERT INTO orders (user_id, total_amount, status, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?)
        RETURNING id;
    `)

	insertOrderItemQueryTest := masterDB.Rebind(`
        INSERT INTO order_items (order_id, book_id, quantity, price, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?);
    `)

	type args struct {
		ctx   context.Context
		order orders.CreateOrderRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *orders.CreateOrderResponse
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "error on begin transaction",
			args: args{
				ctx: context.Background(),
				order: orders.CreateOrderRequest{
					UserID:      1,
					TotalAmount: 100.0,
					Items: []orders.CreateOrderItem{
						{
							BookID:   101,
							Quantity: 2,
							Price:    50.0,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin().WillReturnError(errors.New("failed to begin transaction"))
			},
		},
		{
			name: "error on prepare insert order query",
			args: args{
				ctx: context.Background(),
				order: orders.CreateOrderRequest{
					UserID:      1,
					TotalAmount: 100.0,
					Items: []orders.CreateOrderItem{
						{
							BookID:   101,
							Quantity: 2,
							Price:    50.0,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectPrepare(insertOrderQueryTest).WillReturnError(errors.New("failed to prepare order query"))
				mock.ExpectRollback()
			},
		},
		{
			name: "error on insert order query",
			args: args{
				ctx: context.Background(),
				order: orders.CreateOrderRequest{
					UserID:      1,
					TotalAmount: 100.0,
					Items: []orders.CreateOrderItem{
						{
							BookID:   101,
							Quantity: 2,
							Price:    50.0,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectPrepare(insertOrderQueryTest).ExpectQuery().WillReturnError(errors.New("failed to insert order"))
				mock.ExpectRollback()
			},
		},
		{
			name: "error on prepare insert order item query",
			args: args{
				ctx: context.Background(),
				order: orders.CreateOrderRequest{
					UserID:      1,
					TotalAmount: 100.0,
					Items: []orders.CreateOrderItem{
						{
							BookID:   101,
							Quantity: 2,
							Price:    50.0,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectPrepare(insertOrderQueryTest).ExpectQuery().
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectPrepare(insertOrderItemQueryTest).WillReturnError(errors.New("failed to prepare order item query"))
				mock.ExpectRollback()
			},
		},
		{
			name: "error on insert order item query",
			args: args{
				ctx: context.Background(),
				order: orders.CreateOrderRequest{
					UserID:      1,
					TotalAmount: 100.0,
					Items: []orders.CreateOrderItem{
						{
							BookID:   101,
							Quantity: 2,
							Price:    50.0,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectPrepare(insertOrderQueryTest).ExpectQuery().
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectPrepare(insertOrderItemQueryTest).ExpectExec().WillReturnError(errors.New("failed to insert order item"))
				mock.ExpectRollback()
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				order: orders.CreateOrderRequest{
					UserID:      1,
					TotalAmount: 100.0,
					Items: []orders.CreateOrderItem{
						{
							BookID:   101,
							Quantity: 2,
							Price:    50.0,
						},
					},
				},
			},
			want: &orders.CreateOrderResponse{
				OrderID: 1,
				Status:  orders.OrderStatusNew.String(),
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()
				mock.ExpectPrepare(insertOrderQueryTest).ExpectQuery().
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectPrepare(insertOrderItemQueryTest).ExpectExec().
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &repository{
				masterDB: masterDB,
			}
			got, err := r.InsertOrder(tt.args.ctx, tt.args.order)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repository_GetOrdersByUserID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	slaveDB := internalsql.NewSlaveDB(db, "sqlmock")
	defer func() {
		_ = db.Close()
	}()

	type args struct {
		ctx    context.Context
		userID int64
		limit  int
		offset int
	}
	getOrderQueryTest := slaveDB.Rebind(`
					SELECT id, total_amount, status, created_at, updated_at
					FROM orders
					WHERE user_id = ?
					ORDER BY created_at DESC
					LIMIT ? OFFSET ?
				`)

	getOrderItemQueryTest := slaveDB.Rebind(`
					SELECT id, order_id, book_id, quantity, price
					FROM order_items
					WHERE order_id = ANY(?)
				`)

	tests := []struct {
		name    string
		args    args
		want    []orders.History
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "error when preparing order statement",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				limit:  10,
				offset: 0,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectPrepare(getOrderQueryTest).WillReturnError(errors.New("failed"))
			},
		},
		{
			name: "error when querying order statement",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				limit:  10,
				offset: 0,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectPrepare(getOrderQueryTest).ExpectQuery().
					WithArgs(1, 10, 0).
					WillReturnError(errors.New("failed"))
			},
		},
		{
			name: "error when scanning order row",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				limit:  10,
				offset: 0,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "total_amount", "status", "created_at", "updated_at"}).
					AddRow("invalid_id", 100, "NEW", 1623550814, 1623550814)
				mock.ExpectPrepare(getOrderQueryTest).ExpectQuery().
					WithArgs(1, 10, 0).
					WillReturnRows(rows)
			},
		},
		{
			name: "no orders found for user",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				limit:  10,
				offset: 0,
			},
			want:    nil,
			wantErr: false,
			mockFn: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "total_amount", "status", "created_at", "updated_at"})
				mock.ExpectPrepare(getOrderQueryTest).ExpectQuery().
					WithArgs(1, 10, 0).
					WillReturnRows(rows)
			},
		},
		{
			name: "error when preparing items statement",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				limit:  10,
				offset: 0,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "total_amount", "status", "created_at", "updated_at"}).
					AddRow(1, 100, "NEW", 1623550814, 1623550814)
				mock.ExpectPrepare(getOrderQueryTest).ExpectQuery().
					WithArgs(1, 10, 0).
					WillReturnRows(rows)
				mock.ExpectPrepare(getOrderItemQueryTest).WillReturnError(errors.New("failed"))
			},
		},
		{
			name: "error when querying items statement",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				limit:  10,
				offset: 0,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				rows := sqlmock.NewRows([]string{"id", "total_amount", "status", "created_at", "updated_at"}).
					AddRow(1, 100, "NEW", 1623550814, 1623550814)
				mock.ExpectPrepare(getOrderQueryTest).ExpectQuery().
					WithArgs(1, 10, 0).
					WillReturnRows(rows)
				mock.ExpectPrepare(getOrderItemQueryTest).ExpectQuery().
					WithArgs(pq.Array([]int64{1})).
					WillReturnError(errors.New("failed"))
			},
		},
		{
			name: "error when scanning items row",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				limit:  10,
				offset: 0,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				orderRows := sqlmock.NewRows([]string{"id", "total_amount", "status", "created_at", "updated_at"}).
					AddRow(1, 100, "NEW", 1623550814, 1623550814)
				mock.ExpectPrepare(getOrderQueryTest).ExpectQuery().
					WithArgs(1, 10, 0).
					WillReturnRows(orderRows)

				itemsRows := sqlmock.NewRows([]string{"id", "order_id", "book_id", "quantity", "price"}).
					AddRow("invalid_id", 1, 1, 2, 50)
				mock.ExpectPrepare(getOrderItemQueryTest).ExpectQuery().
					WithArgs(pq.Array([]int64{1})).
					WillReturnRows(itemsRows)
			},
		},
		{
			name: "successful retrieval of order history with items",
			args: args{
				ctx:    context.Background(),
				userID: 1,
				limit:  10,
				offset: 0,
			},
			want: []orders.History{
				{
					ID:          1,
					TotalAmount: 100,
					Status:      "NEW",
					CreatedAt:   1623550814,
					UpdatedAt:   1623550814,
					Items: []orders.ItemHistory{
						{
							ID:       1,
							BookID:   1,
							Quantity: 2,
							Price:    50,
						},
					},
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				orderRows := sqlmock.NewRows([]string{"id", "total_amount", "status", "created_at", "updated_at"}).
					AddRow(1, 100, "NEW", 1623550814, 1623550814)
				mock.ExpectPrepare(getOrderQueryTest).ExpectQuery().
					WithArgs(1, 10, 0).
					WillReturnRows(orderRows)

				itemsRows := sqlmock.NewRows([]string{"id", "order_id", "book_id", "quantity", "price"}).
					AddRow(1, 1, 1, 2, 50)
				mock.ExpectPrepare(getOrderItemQueryTest).ExpectQuery().
					WithArgs(pq.Array([]int64{1})).
					WillReturnRows(itemsRows)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &repository{
				slaveDB: slaveDB,
			}
			got, err := r.GetOrdersByUserID(tt.args.ctx, tt.args.userID, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrdersByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrdersByUserID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
