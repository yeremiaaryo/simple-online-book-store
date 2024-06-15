package orders

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/orders"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestHandler_CreateOrder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOrdersUC := NewMockordersUsecase(mockCtrl)

	type args struct {
		payload string
		userID  int64
	}
	tests := []struct {
		name   string
		args   args
		want   string
		mockFn func(args args)
	}{
		{
			name: "error invalid user id",
			args: args{
				payload: `invalid_json`,
			},
			want: `{"error":"userID not found", "order_id":0, "result":false, "status":""}`,
			mockFn: func(args args) {

			},
		},
		{
			name: "error bind",
			args: args{
				payload: `invalid_json`,
				userID:  1,
			},
			want: `{"error":"code=400, message=Syntax error: offset=1, error=invalid character 'i' looking for beginning of value, internal=invalid character 'i' looking for beginning of value", "order_id":0, "result":false, "status":""}`,
			mockFn: func(args args) {

			},
		},
		{
			name: "error validate",
			args: args{
				payload: `{}`,
				userID:  1,
			},
			want: `{"result":false,"error":"Key: 'CreateOrderRequest.TotalAmount' Error:Field validation for 'TotalAmount' failed on the 'required' tag\nKey: 'CreateOrderRequest.Items' Error:Field validation for 'Items' failed on the 'required' tag", "order_id":0, "result":false, "status":""}`,
			mockFn: func(args args) {

			},
		},
		{
			name: "error InsertOrder",
			args: args{
				payload: `{"items":[{"book_id":1,"quantity":2,"price":50.0}],"total_amount":100.0}`,
				userID:  1,
			},
			want: `{"error":"failed to insert order", "order_id":0, "result":false, "status":""}`,
			mockFn: func(args args) {
				mockOrdersUC.EXPECT().InsertOrder(gomock.Any(), gomock.Any()).Return(nil, errors.New("failed to insert order"))
			},
		},
		{
			name: "success",
			args: args{
				payload: `{"items":[{"book_id":1,"quantity":2,"price":50.0}],"total_amount":100.0}`,
				userID:  1,
			},
			want: `{"order_id":1, "result":true, "status":"NEW"}`,
			mockFn: func(args args) {
				mockOrdersUC.EXPECT().InsertOrder(gomock.Any(), gomock.Any()).Return(&orders.CreateOrderResponse{
					OrderID: 1,
					Status:  orders.OrderStatusNew.String(),
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			h := &Handler{
				ordersUsecase: mockOrdersUC,
			}
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/order", strings.NewReader(tt.args.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if tt.args.userID != 0 {
				c.Set("userID", tt.args.userID)
			}
			if assert.NoError(t, h.CreateOrder(c)) {
				assert.JSONEq(t, tt.want, rec.Body.String())
			}
		})
	}
}

func TestHandler_GetOrderHistory(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOrdersUC := NewMockordersUsecase(mockCtrl)

	type args struct {
		userID    int64
		pageIndex string
		pageSize  string
	}

	tests := []struct {
		name   string
		args   args
		want   string
		mockFn func(userID int64, pageIndex, pageSize string)
	}{
		{
			name: "error invalid user id",
			args: args{
				userID:    0,
				pageIndex: "1",
				pageSize:  "10",
			},
			want:   `{"data":null, "error":"userID not found", "result":false}`,
			mockFn: func(userID int64, pageIndex, pageSize string) {},
		},
		{
			name: "error GetOrdersByUserID",
			args: args{
				userID:    1,
				pageIndex: "invalid, use default",
				pageSize:  "invalid, use default",
			},
			want: `{"data":null, "error":"failed to retrieve order history", "result":false}`,
			mockFn: func(userID int64, pageIndex, pageSize string) {
				mockOrdersUC.EXPECT().GetOrdersByUserID(gomock.Any(), userID, 1, 10).Return(nil, errors.New("failed to retrieve order history"))
			},
		},
		{
			name: "success",
			args: args{
				userID:    1,
				pageIndex: "1",
				pageSize:  "10",
			},
			want: `{"data":[{"order_id":1,"total_amount":100.0,"status":"NEW","created_at":1623800000,"updated_at":1623800000,"items":[{"item_id":1,"book_id":1,"quantity":2,"price":50.0}]}], "result":true}`,
			mockFn: func(userID int64, pageIndex, pageSize string) {
				mockOrdersUC.EXPECT().GetOrdersByUserID(gomock.Any(), userID, 1, 10).Return([]orders.History{
					{
						ID:          1,
						TotalAmount: 100.0,
						Status:      "NEW",
						CreatedAt:   1623800000,
						UpdatedAt:   1623800000,
						Items: []orders.ItemHistory{
							{
								ID:       1,
								BookID:   1,
								Quantity: 2,
								Price:    50.0,
							},
						},
					},
				}, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args.userID, tt.args.pageIndex, tt.args.pageSize)

			h := &Handler{
				ordersUsecase: mockOrdersUC,
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/order", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			if tt.args.userID != 0 {
				c.Set("userID", tt.args.userID)
			}
			c.QueryParams().Add("page_index", tt.args.pageIndex)
			c.QueryParams().Add("page_size", tt.args.pageSize)

			if assert.NoError(t, h.GetOrderHistory(c)) {
				assert.JSONEq(t, tt.want, rec.Body.String())
			}
		})
	}
}
