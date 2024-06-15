package internalsql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"
	sqlxtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
)

type DB interface {
	Rebind(query string) string
	Ping() error
	Close() error
}

type Statement interface {
	Close() error
	QueryRowxContext(ctx context.Context, args ...interface{}) *sqlx.Row
	QueryxContext(ctx context.Context, args ...interface{}) (*sqlx.Rows, error)
}

type MasterDB interface {
	DB
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
	PreparexContext(ctx context.Context, query string) (MasterStatement, error)
}

type SlaveDB interface {
	DB
	PreparexContext(ctx context.Context, query string) (SlaveStatement, error)
}

type MasterStatement interface {
	Statement
	ExecContext(ctx context.Context, args ...interface{}) (sql.Result, error)
}

type SlaveStatement interface {
	Statement
	GetContext(ctx context.Context, dest interface{}, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, args ...interface{}) error
}

type slaveDBImpl struct {
	*sqlx.DB
}

type masterDBImpl struct {
	*sqlx.DB
}

func OpenMasterDB(driver, dsn string) (MasterDB, error) {
	sqltrace.Register("postgres", &pq.Driver{}, sqltrace.WithServiceName("gotu"))
	db, err := sqlxtrace.Open(driver, dsn)

	if err != nil {
		return nil, err
	}

	return &masterDBImpl{
		DB: db,
	}, db.Ping()
}
func (mdb *masterDBImpl) PreparexContext(ctx context.Context, query string) (MasterStatement, error) {
	return mdb.DB.PreparexContext(ctx, query)
}

func OpenSlaveDB(driver, dsn string) (SlaveDB, error) {
	sqltrace.Register("postgres", &pq.Driver{}, sqltrace.WithServiceName("gotu"))
	db, err := sqlxtrace.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	return &slaveDBImpl{
		DB: db,
	}, db.Ping()
}

func (sd *slaveDBImpl) PreparexContext(ctx context.Context, query string) (SlaveStatement, error) {
	return sd.DB.PreparexContext(ctx, query)
}

// NewMasterDB creates new MasterDB object from existing sql.DB object.
func NewMasterDB(db *sql.DB, driverName string) MasterDB {
	return &masterDBImpl{
		DB: sqlx.NewDb(db, driverName),
	}
}

// NewSlaveDB creates new SlaveDB object from existing sql.DB object.
func NewSlaveDB(db *sql.DB, driverName string) SlaveDB {
	return &slaveDBImpl{
		DB: sqlx.NewDb(db, driverName),
	}
}
