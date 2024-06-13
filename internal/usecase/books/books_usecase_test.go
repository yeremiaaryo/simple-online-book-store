package books

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/books"
	"reflect"
	"testing"
)

func Test_usecase_GetBooks(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockBooksRepo := NewMockbooksRepository(mockCtrl)
	type args struct {
		ctx       context.Context
		search    string
		pageSize  int
		pageIndex int
	}
	tests := []struct {
		name    string
		args    args
		want    []books.Model
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "error",
			args: args{
				ctx:       context.Background(),
				search:    "Book",
				pageSize:  10,
				pageIndex: 1,
			},

			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mockBooksRepo.EXPECT().GetBooks(args.ctx, args.search, args.pageSize, 0).Return(nil, errors.New("failed"))
			},
		},
		{
			name: "valid search, page 1, page size 10",
			args: args{
				ctx:       context.Background(),
				search:    "Book",
				pageSize:  10,
				pageIndex: 1,
			},

			want: []books.Model{
				{ID: 1, Title: "Book 1", Author: "Author 1", ISBN: "123456789", CreatedAt: 1623582000, UpdatedAt: 1623582000},
				{ID: 2, Title: "Book 2", Author: "Author 2", ISBN: "987654321", CreatedAt: 1623582000, UpdatedAt: 1623582000},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockBooksRepo.EXPECT().GetBooks(args.ctx, args.search, args.pageSize, 0).Return([]books.Model{
					{ID: 1, Title: "Book 1", Author: "Author 1", ISBN: "123456789", CreatedAt: 1623582000, UpdatedAt: 1623582000},
					{ID: 2, Title: "Book 2", Author: "Author 2", ISBN: "987654321", CreatedAt: 1623582000, UpdatedAt: 1623582000},
				}, nil)
			},
		},
		{
			name: "invalid page index (negative)",
			args: args{
				ctx:       context.Background(),
				search:    "",
				pageSize:  10,
				pageIndex: -1,
			},
			want: []books.Model{
				{ID: 1, Title: "Book 1", Author: "Author 1", ISBN: "123456789", CreatedAt: 1623582000, UpdatedAt: 1623582000},
				{ID: 2, Title: "Book 2", Author: "Author 2", ISBN: "987654321", CreatedAt: 1623582000, UpdatedAt: 1623582000},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockBooksRepo.EXPECT().GetBooks(args.ctx, args.search, args.pageSize, 0).Return([]books.Model{
					{ID: 1, Title: "Book 1", Author: "Author 1", ISBN: "123456789", CreatedAt: 1623582000, UpdatedAt: 1623582000},
					{ID: 2, Title: "Book 2", Author: "Author 2", ISBN: "987654321", CreatedAt: 1623582000, UpdatedAt: 1623582000},
				}, nil)
			},
		},
		{
			name: "invalid page size (zero)",
			args: args{
				ctx:       context.Background(),
				search:    "",
				pageSize:  0,
				pageIndex: 1,
			},
			want: []books.Model{
				{ID: 1, Title: "Book 1", Author: "Author 1", ISBN: "123456789", CreatedAt: 1623582000, UpdatedAt: 1623582000},
				{ID: 2, Title: "Book 2", Author: "Author 2", ISBN: "987654321", CreatedAt: 1623582000, UpdatedAt: 1623582000},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockBooksRepo.EXPECT().GetBooks(args.ctx, args.search, 10, 0).Return([]books.Model{
					{ID: 1, Title: "Book 1", Author: "Author 1", ISBN: "123456789", CreatedAt: 1623582000, UpdatedAt: 1623582000},
					{ID: 2, Title: "Book 2", Author: "Author 2", ISBN: "987654321", CreatedAt: 1623582000, UpdatedAt: 1623582000},
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			u := &usecase{
				booksRepository: mockBooksRepo,
			}
			got, err := u.GetBooks(tt.args.ctx, tt.args.search, tt.args.pageSize, tt.args.pageIndex)
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
