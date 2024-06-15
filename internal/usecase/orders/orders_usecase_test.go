package orders

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/books"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/orders"
	"reflect"
	"testing"
)

func Test_usecase_InsertOrder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBooksRepo := NewMockbooksRepository(mockCtrl)
	mockOrdersRepo := NewMockordersRepository(mockCtrl)

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
			name: "error when getting book by IDs",
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
				mockBooksRepo.EXPECT().GetBookByIDs(args.ctx, gomock.Any()).Return(nil, errors.New("failed to get books by IDs"))
			},
		},
		{
			name: "error due to different book price",
			args: args{
				ctx: context.Background(),
				order: orders.CreateOrderRequest{
					UserID:      1,
					TotalAmount: 100.0,
					Items: []orders.CreateOrderItem{
						{
							BookID:   101,
							Quantity: 2,
							Price:    60.0,
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockBooksRepo.EXPECT().GetBookByIDs(args.ctx, gomock.Any()).Return(map[int64]books.Model{
					101: {
						ID:    101,
						Title: "Book 101",
						Price: 50.0,
					},
				}, nil)
			},
		},
		{
			name: "error due to total amount mismatch",
			args: args{
				ctx: context.Background(),
				order: orders.CreateOrderRequest{
					UserID:      1,
					TotalAmount: 200.0, // Incorrect total amount
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
				mockBooksRepo.EXPECT().GetBookByIDs(args.ctx, gomock.Any()).Return(map[int64]books.Model{
					101: {
						ID:    101,
						Title: "Book 101",
						Price: 50.0,
					},
				}, nil)
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
				mockBooksRepo.EXPECT().GetBookByIDs(args.ctx, gomock.Any()).Return(map[int64]books.Model{
					101: {
						ID:    101,
						Title: "Book 101",
						Price: 50.0,
					},
				}, nil)
				mockOrdersRepo.EXPECT().InsertOrder(args.ctx, args.order).Return(&orders.CreateOrderResponse{
					OrderID: 1,
					Status:  orders.OrderStatusNew.String(),
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			u := &usecase{
				booksRepository:  mockBooksRepo,
				ordersRepository: mockOrdersRepo,
			}
			got, err := u.InsertOrder(tt.args.ctx, tt.args.order)
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

func Test_usecase_GetOrdersByUserID(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOrdersRepo := NewMockordersRepository(mockCtrl)

	type args struct {
		ctx       context.Context
		userID    int64
		pageIndex int
		pageSize  int
	}
	tests := []struct {
		name    string
		args    args
		want    []orders.History
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "error in repository",
			args: args{
				ctx:       context.Background(),
				userID:    1,
				pageIndex: 1,
				pageSize:  10,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockOrdersRepo.EXPECT().GetOrdersByUserID(args.ctx, args.userID, 10, 0).Return(nil, errors.New("repository error"))
			},
		},
		{
			name: "successful retrieval",
			args: args{
				ctx:       context.Background(),
				userID:    1,
				pageIndex: 1,
				pageSize:  10,
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
							Price:    10,
						},
					},
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockOrdersRepo.EXPECT().GetOrdersByUserID(args.ctx, args.userID, 10, 0).Return([]orders.History{
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
								Price:    10,
							},
						},
					},
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			u := &usecase{
				ordersRepository: mockOrdersRepo,
			}
			got, err := u.GetOrdersByUserID(tt.args.ctx, tt.args.userID, tt.args.pageIndex, tt.args.pageSize)
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
