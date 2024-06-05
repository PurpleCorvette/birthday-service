package db

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type DB interface {
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Query(ctx context.Context, sql string, args ...interface{}) (Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Close()
}

type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close()
}

type PgxDB struct {
	Pool *pgxpool.Pool
}

func (db *PgxDB) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return db.Pool.QueryRow(ctx, sql, args...)
}

func (db *PgxDB) Query(ctx context.Context, sql string, args ...interface{}) (Rows, error) {
	rows, err := db.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return &PgxRows{rows}, nil
}

func (db *PgxDB) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return db.Pool.Exec(ctx, sql, args...)
}

func (db *PgxDB) Close() {
	db.Pool.Close()
}

type PgxRows struct {
	pgx.Rows
}

func (r *PgxRows) Next() bool {
	return r.Rows.Next()
}

func (r *PgxRows) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r *PgxRows) Close() {
	r.Rows.Close()
}

func ConnectDatabase(url string) (DB, error) {
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return &PgxDB{Pool: pool}, nil
}
