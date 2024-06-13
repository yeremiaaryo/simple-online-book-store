package users

import (
	"context"
	"database/sql"
	"github.com/yeremiaaryo/gotu-assignment/internal/model/users"
	"github.com/yeremiaaryo/gotu-assignment/pkg/internalsql"
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

func (r *repository) GetUser(ctx context.Context, email string) (*users.Model, error) {
	query := getUsersQuery + ` WHERE email = ?`
	rebindQuery := r.slaveDB.Rebind(query)

	stmt, err := r.slaveDB.PreparexContext(ctx, rebindQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user users.Model
	err = stmt.GetContext(ctx, &user, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) InsertUser(ctx context.Context, model users.Model) (*users.Model, error) {
	rebindQuery := r.masterDB.Rebind(insertUserQuery)

	stmt, err := r.masterDB.PreparexContext(ctx, rebindQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRowxContext(ctx, model.Email, model.Password, model.CreatedAt, model.UpdatedAt).Scan(&model.ID)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
