package users

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/users"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func TestHandler_CreateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUsersUC := NewMockusersUsecase(mockCtrl)
	type args struct {
		payload string
	}
	tests := []struct {
		name   string
		args   args
		want   string
		mockFn func(args args)
	}{
		{
			name: "error bind",
			args: args{
				payload: `failed`,
			},
			want: `{"result":false,"error":"code=400, message=Syntax error: offset=3, error=invalid character 'i' in literal false (expecting 'l'), internal=invalid character 'i' in literal false (expecting 'l')","user":null}`,
			mockFn: func(args args) {

			},
		},
		{
			name: "error validate",
			args: args{
				payload: `{}`,
			},
			want: `{"result":false,"error":"Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'CreateUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag","user":null}`,
			mockFn: func(args args) {

			},
		},
		{
			name: "error CreateUser",
			args: args{
				payload: `{"email":"email@email.com","password":"password"}`,
			},
			want: `{"result":false,"error":"email already exists","user":null}`,
			mockFn: func(args args) {
				mockUsersUC.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, errors.New("email already exists"))
			},
		},
		{
			name: "success",
			args: args{
				payload: `{"email":"email@email.com","password":"password"}`,
			},
			want: `{"result":true,"user":{"id":1,"email":"email@email.com","created_at":1714580787000,"updated_at":1714580787000}}`,
			mockFn: func(args args) {
				mockUsersUC.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(&users.Model{
					ID:        1,
					Email:     "email@email.com",
					CreatedAt: 1714580787000,
					UpdatedAt: 1714580787000,
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			h := &Handler{
				usersUsecase: mockUsersUC,
			}
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tt.args.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if assert.NoError(t, h.CreateUser(c)) {
				assert.JSONEq(t, tt.want, rec.Body.String())
			}
		})
	}
}

func TestHandler_Login(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUsersUC := NewMockusersUsecase(mockCtrl)

	type args struct {
		payload string
	}
	tests := []struct {
		name   string
		args   args
		want   string
		mockFn func(args args)
	}{
		{
			name: "error bind",
			args: args{
				payload: `failed`,
			},
			want: `{"result":false,"error":"code=400, message=Syntax error: offset=3, error=invalid character 'i' in literal false (expecting 'l'), internal=invalid character 'i' in literal false (expecting 'l')","token":""}`,
			mockFn: func(args args) {

			},
		},
		{
			name: "error validate",
			args: args{
				payload: `{}`,
			},
			want: `{"result":false,"error":"Key: 'LoginRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'LoginRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag","token":""}`,
			mockFn: func(args args) {

			},
		},
		{
			name: "error Login",
			args: args{
				payload: `{"email":"email@email.com","password":"12345"}`,
			},
			want: `{"result":false,"error":"invalid email or password","token":""}`,
			mockFn: func(args args) {
				mockUsersUC.EXPECT().Login(gomock.Any(), gomock.Any()).Return("", errors.New("invalid email or password"))
			},
		},
		{
			name: "success",
			args: args{
				payload: `{"email":"email@email.com","password":"12345"}`,
			},
			want: `{"result":true,"token":"token"}`,
			mockFn: func(args args) {
				mockUsersUC.EXPECT().Login(gomock.Any(), gomock.Any()).Return("token", nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			h := &Handler{
				usersUsecase: mockUsersUC,
			}
			e := echo.New()
			e.Validator = &CustomValidator{validator: validator.New()}
			req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(tt.args.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			if assert.NoError(t, h.Login(c)) {
				fmt.Println(tt.name, rec.Body.String())
				assert.JSONEq(t, tt.want, rec.Body.String())
			}
		})
	}
}
