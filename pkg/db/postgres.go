package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, connString string) (*Postgres, error) {
	if pgInstance != nil {
		return pgInstance, nil
	}
	var (
		err    error
		dbPool *pgxpool.Pool
	)
	pgOnce.Do(func() {
		dbPool, err = pgxpool.New(ctx, connString)
		pgInstance = &Postgres{Pool: dbPool}
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	return pgInstance, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.Pool.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.Pool.Close()
}
