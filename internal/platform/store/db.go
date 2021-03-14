package store

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func Crdb(ctx context.Context, url string) (*sqlx.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	return sqlx.NewDb(db, "pq"), nil
}
