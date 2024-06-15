package users

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/yeremiaaryo/gotu-assignment/internal/configs"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/users"
)

func Test_usecase_CreateUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUsersRepo := NewMockusersRepository(mockCtrl)

	type args struct {
		ctx context.Context
		req users.CreateUserRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *users.Model
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "error when GetUser",
			args: args{
				ctx: context.Background(),
				req: users.CreateUserRequest{
					Email: "email@email.com",
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockUsersRepo.EXPECT().GetUser(gomock.Any(), args.req.Email).Return(nil, errors.New("failed"))
			},
		},
		{
			name: "error when GetUser, user exist",
			args: args{
				ctx: context.Background(),
				req: users.CreateUserRequest{
					Email: "email@email.com",
				},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockUsersRepo.EXPECT().GetUser(gomock.Any(), args.req.Email).Return(&users.Model{}, nil)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: users.CreateUserRequest{
					Email:    "email@email.com",
					Password: "pass",
				},
			},
			want: &users.Model{
				ID:    1,
				Email: "email@email.com",
			},
			wantErr: false,
			mockFn: func(args args) {
				mockUsersRepo.EXPECT().GetUser(gomock.Any(), args.req.Email).Return(nil, nil)
				mockUsersRepo.EXPECT().InsertUser(gomock.Any(), gomock.Any()).Return(&users.Model{
					ID:    1,
					Email: "email@email.com",
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			u := &usecase{
				usersRepository: mockUsersRepo,
			}
			got, err := u.CreateUser(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_usecase_Login(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUsersRepo := NewMockusersRepository(mockCtrl)
	mockRedis := NewMockredis(mockCtrl)

	type args struct {
		ctx context.Context
		req users.LoginRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "error when GetUser",
			args: args{
				ctx: context.Background(),
				req: users.LoginRequest{
					Email: "email@email.com",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockUsersRepo.EXPECT().GetUser(gomock.Any(), args.req.Email).Return(nil, errors.New("failed"))
			},
		},
		{
			name: "error when GetUser, user not found",
			args: args{
				ctx: context.Background(),
				req: users.LoginRequest{
					Email: "email@email.com",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockUsersRepo.EXPECT().GetUser(gomock.Any(), args.req.Email).Return(nil, nil)
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: users.LoginRequest{
					Email:    "email@email.com",
					Password: "12345",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockUsersRepo.EXPECT().GetUser(gomock.Any(), args.req.Email).Return(&users.Model{ID: 123, Password: `$2a$10$Kes/ccWjDAw01VM1STV8mePua4YOpMwldDqlLq7GltRvJr/zdj7zq`}, nil)
				mockRedis.EXPECT().Set("token:123", gomock.Any(), int64((24 * time.Hour).Seconds()))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			u := &usecase{
				usersRepository: mockUsersRepo,
				redis:           mockRedis,
				cfg: &configs.Config{
					Service: configs.Service{
						SecretKey: "secretkey",
					},
				},
			}
			got, err := u.Login(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" && !tt.wantErr {
				t.Errorf("Login() got = %v", got)
			}
		})
	}
}
