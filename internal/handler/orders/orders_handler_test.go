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
			req := httptest.NewRequest(http.MethodPost, "/create-order", strings.NewReader(tt.args.payload))
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
