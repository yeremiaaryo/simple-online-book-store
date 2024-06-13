package books

import (
	"context"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/books"
	"github.com/yeremiaaryo/gotu-assignment/pkg/internalsql"
	"strings"
)

type repository struct {
	masterDB internalsql.MasterDB
	slaveDB  internalsql.SlaveDB
}

func New(masterDB internalsql.MasterDB, slaveDB internalsql.SlaveDB) *repository {
	r := repository{
		masterDB: masterDB,
		slaveDB:  slaveDB,
	}

	return &r
}

func (r *repository) GetBooks(ctx context.Context, search string, limit, offset int) ([]books.Model, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString(queryGetBooks)

	var args []interface{}

	if search != "" {
		queryBuilder.WriteString(` WHERE lower(title) ILIKE lower(?) OR lower(author) ILIKE lower(?)`)
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	queryBuilder.WriteString(` LIMIT ? OFFSET ?`)
	args = append(args, limit, offset)

	query := queryBuilder.String()
	rebindQuery := r.slaveDB.Rebind(query)

	stmt, err := r.slaveDB.PreparexContext(ctx, rebindQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var bookList []books.Model
	err = stmt.SelectContext(ctx, &bookList, args...)
	if err != nil {
		return nil, err
	}
	return bookList, nil
}
