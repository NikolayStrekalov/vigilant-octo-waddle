package pgstorage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func Ping(dsn string) bool {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return false
	}
	defer func() {
		_ = conn.Close(context.Background())
	}()
	return true
}
