package orders

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
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
