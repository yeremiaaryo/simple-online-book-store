package users

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/users"
	"github.com/yeremiaaryo/gotu-assignment/pkg/internalsql"
)

func Test_repository_GetUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	masterDB := internalsql.NewMasterDB(db, "sqlmock")
	slaveDB := internalsql.NewSlaveDB(db, "sqlmock")
	type args struct {
		ctx   context.Context
		email string
	}
	defer func() {
		_ = db.Close()
	}()

	query := `SELECT 
				id, email, password, created_at, updated_at 
			FROM 
				users WHERE email = ? `
	rebindQuery := slaveDB.Rebind(query)

	tests := []struct {
		name    string
		args    args
		want    *users.Model
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "error when prepare context",
			args: args{
				ctx:   context.Background(),
				email: "email@email.com",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectPrepare(rebindQuery).WillReturnError(errors.New("failed"))
			},
		},
		{
			name: "error when get context",
			args: args{
				ctx:   context.Background(),
				email: "email@email.com",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectPrepare(rebindQuery).ExpectQuery().WillReturnError(errors.New("failed"))
			},
		},
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				email: "email@email.com",
			},
			want: &users.Model{
				ID:        1,
				Email:     "email@email.com",
				Password:  "password",
				CreatedAt: 1714641784000,
				UpdatedAt: 1714641784000,
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectPrepare(rebindQuery).ExpectQuery().
					WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "created_at", "updated_at"}).
						AddRow(1, "email@email.com", "password", 1714641784000, 1714641784000))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &repository{
				masterDB: masterDB,
				slaveDB:  slaveDB,
			}
			got, err := r.GetUser(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repository_InsertUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	masterDB := internalsql.NewMasterDB(db, "sqlmock")
	slaveDB := internalsql.NewSlaveDB(db, "sqlmock")
	defer func() {
		_ = db.Close()
	}()

	query := `INSERT INTO users
							(email, password, created_at, updated_at)
							VALUES(?, ?, ?, ?) RETURNING id;`
	rebindQuery := masterDB.Rebind(query)

	type args struct {
		ctx   context.Context
		model users.Model
	}
	tests := []struct {
		name    string
		args    args
		want    *users.Model
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "error when prepare context",
			args: args{
				ctx:   context.Background(),
				model: users.Model{},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectPrepare(rebindQuery).WillReturnError(errors.New("failed"))
			},
		},
		{
			name: "error when get context",
			args: args{
				ctx:   context.Background(),
				model: users.Model{},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectPrepare(rebindQuery).ExpectQuery().WillReturnError(errors.New("failed"))
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				model: users.Model{
					Email:     "email@email.com",
					Password:  "password",
					CreatedAt: 1714641784000,
					UpdatedAt: 1714641784000,
				},
			},
			want: &users.Model{
				ID:        1,
				Email:     "email@email.com",
				Password:  "password",
				CreatedAt: 1714641784000,
				UpdatedAt: 1714641784000,
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectPrepare(rebindQuery).ExpectQuery().
					WillReturnRows(sqlmock.NewRows([]string{"id"}).
						AddRow(1))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := &repository{
				masterDB: masterDB,
				slaveDB:  slaveDB,
			}
			got, err := r.InsertUser(tt.args.ctx, tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("InsertUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InsertUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
