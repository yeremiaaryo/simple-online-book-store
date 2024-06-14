package books

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/books"
	"github.com/yeremiaaryo/gotu-assignment/pkg/internalsql"
	"reflect"
	"testing"
)

func Test_repository_GetBooks(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	masterDB := internalsql.NewMasterDB(db, "sqlmock")
	slaveDB := internalsql.NewSlaveDB(db, "sqlmock")
	defer func() {
		_ = db.Close()
	}()

	type args struct {
		ctx    context.Context
		search string
		limit  int
		offset int
	}
	tests := []struct {
		name    string
		args    args
		want    []books.Model
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "error when preparing statement",
			args: args{
				ctx:    context.Background(),
				search: "Orwell",
				limit:  10,
				offset: 0,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectPrepare(`SELECT id, title, author, isbn, published_date,  price, created_at, updated_at FROM books WHERE lower(title) ILIKE lower(?) OR lower(author) ILIKE lower(?) LIMIT ? OFFSET ?`).
					WillReturnError(errors.New("failed"))
			},
		},
		{
			name: "error when querying statement",
			args: args{
				ctx:    context.Background(),
				search: "Orwell",
				limit:  10,
				offset: 0,
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectPrepare(`SELECT id, title, author, isbn, published_date, price, created_at, updated_at FROM books WHERE lower(title) ILIKE lower(?) OR lower(author) ILIKE lower(?) LIMIT ? OFFSET ?`).
					ExpectQuery().
					WithArgs("%Orwell%", "%Orwell%", 10, 0).
					WillReturnError(errors.New("failed"))
			},
		},
		{
			name: "search by title success",
			args: args{
				ctx:    context.Background(),
				search: "Orwell",
				limit:  10,
				offset: 0,
			},
			want: []books.Model{
				{
					ID:        1,
					Title:     "1984",
					Author:    "George Orwell",
					ISBN:      "9780451524935",
					Price:     9.99,
					CreatedAt: 1623550814,
					UpdatedAt: 1623550814,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectPrepare(`SELECT id, title, author, isbn, published_date, price, created_at, updated_at FROM books WHERE lower(title) ILIKE lower(?) OR lower(author) ILIKE lower(?) LIMIT ? OFFSET ?`).
					ExpectQuery().
					WithArgs("%Orwell%", "%Orwell%", 10, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "isbn", "price", "created_at", "updated_at"}).
						AddRow(1, "1984", "George Orwell", "9780451524935", 9.99, 1623550814, 1623550814))
			},
		},
		{
			name: "no search term",
			args: args{
				ctx:    context.Background(),
				search: "",
				limit:  10,
				offset: 0,
			},
			want: []books.Model{
				{
					ID:        1,
					Title:     "1984",
					Author:    "George Orwell",
					ISBN:      "9780451524935",
					Price:     9.99,
					CreatedAt: 1623550814,
					UpdatedAt: 1623550814,
				},
				{
					ID:        2,
					Title:     "Animal Farm",
					Author:    "George Orwell",
					ISBN:      "9780451526342",
					Price:     8.99,
					CreatedAt: 1623550814,
					UpdatedAt: 1623550814,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectPrepare(`SELECT id, title, author, isbn, published_date, price, created_at, updated_at FROM books LIMIT ? OFFSET ?`).
					ExpectQuery().
					WithArgs(10, 0).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "isbn", "price", "created_at", "updated_at"}).
						AddRow(1, "1984", "George Orwell", "9780451524935", 9.99, 1623550814, 1623550814).
						AddRow(2, "Animal Farm", "George Orwell", "9780451526342", 8.99, 1623550814, 1623550814))
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
			got, err := r.GetBooks(tt.args.ctx, tt.args.search, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBooks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_repository_GetBookByIDs(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	masterDB := internalsql.NewMasterDB(db, "sqlmock")
	slaveDB := internalsql.NewSlaveDB(db, "sqlmock")
	defer func() {
		_ = db.Close()
	}()

	type args struct {
		ctx context.Context
		ids []int64
	}
	selectQuery := masterDB.Rebind(`SELECT id, title, author, isbn, published_date, price, created_at, updated_at FROM books WHERE id = ANY(?)`)
	tests := []struct {
		name    string
		args    args
		want    map[int64]books.Model
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "error preparing statement",
			args: args{
				ctx: context.Background(),
				ids: []int64{1, 2},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectPrepare(selectQuery).
					WillReturnError(errors.New("failed"))
			},
		},
		{
			name: "error querying statement",
			args: args{
				ctx: context.Background(),
				ids: []int64{1, 2},
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectPrepare(selectQuery).
					ExpectQuery().
					WithArgs(pq.Array(args.ids)).
					WillReturnError(errors.New("failed"))
			},
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				ids: []int64{1, 2},
			},
			want: map[int64]books.Model{
				1: {
					ID:        1,
					Title:     "1984",
					Author:    "George Orwell",
					ISBN:      "9780451524935",
					Price:     9.99,
					CreatedAt: 1623550814,
					UpdatedAt: 1623550814,
				},
				2: {
					ID:        2,
					Title:     "Animal Farm",
					Author:    "George Orwell",
					ISBN:      "9780451526342",
					Price:     8.99,
					CreatedAt: 1623550814,
					UpdatedAt: 1623550814,
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectPrepare(selectQuery).
					ExpectQuery().
					WithArgs(pq.Array(args.ids)).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "isbn", "price", "created_at", "updated_at"}).
						AddRow(1, "1984", "George Orwell", "9780451524935", 9.99, 1623550814, 1623550814).
						AddRow(2, "Animal Farm", "George Orwell", "9780451526342", 8.99, 1623550814, 1623550814))
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
			got, err := r.GetBookByIDs(tt.args.ctx, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBookByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBookByIDs() got = %v, want %v", got, tt.want)
			}
		})
	}
}
