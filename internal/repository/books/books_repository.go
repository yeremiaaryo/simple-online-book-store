package books

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/lib/pq"
	"github.com/yeremiaaryo/gotu-assignment/internal/constant"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/books"
	"github.com/yeremiaaryo/gotu-assignment/pkg/internalsql"
	"strings"
	"time"
)

//go:generate mockgen -package=books -source=books_repository.go -destination=books_repository_mock_test.go
type redis interface {
	Get(key string, field ...interface{}) (string, error)
	Set(key string, value string, ttl int64, field ...interface{}) (interface{}, error)
}

type repository struct {
	masterDB internalsql.MasterDB
	slaveDB  internalsql.SlaveDB
	redis    redis
}

func New(masterDB internalsql.MasterDB, slaveDB internalsql.SlaveDB, redis redis) *repository {
	r := repository{
		masterDB: masterDB,
		slaveDB:  slaveDB,
		redis:    redis,
	}

	return &r
}

func (r *repository) GetBooks(ctx context.Context, search string, limit, offset int) ([]books.Model, error) {
	var bookList []books.Model

	redisKey := fmt.Sprintf(constant.RedisKeyBooks, search, limit, offset)
	resStr, err := r.redis.Get(redisKey)
	if err == nil && resStr != "" {
		err = jsoniter.Unmarshal([]byte(resStr), &bookList)
		if err == nil {
			return bookList, nil
		}
	}

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

	err = stmt.SelectContext(ctx, &bookList, args...)
	if err != nil {
		return nil, err
	}

	val, err := jsoniter.MarshalToString(bookList)
	if err != nil {
		return bookList, nil // still return no error, just error on set redis shouldn't block user journey
	}
	fmt.Println(val)
	_, _ = r.redis.Set(redisKey, val, int64((30 * time.Second).Seconds()))
	return bookList, nil
}

func (r *repository) GetBookByIDs(ctx context.Context, ids []int64) (map[int64]books.Model, error) {
	var queryBuilder strings.Builder
	queryBuilder.WriteString(queryGetBooks)

	var args []interface{}

	queryBuilder.WriteString(` WHERE id = ANY(?)`)
	args = append(args, pq.Array(ids))

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

	result := make(map[int64]books.Model, 0)
	for _, book := range bookList {
		result[book.ID] = book
	}
	return result, nil
}
