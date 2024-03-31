package tools

import (
	"context"
	"time"

	"github.com/L1LSunflower/lina/pkg/db"
)

const dependsKey = "depends"

type depends struct {
	pg *db.Postgres
}

func CtxWithDepends(parentCtx context.Context, db *db.Postgres) context.Context {
	return context.WithValue(parentCtx, dependsKey, &depends{pg: db})
}

func CtxWithTimeout(parentCtx context.Context, timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(parentCtx, time.Duration(timeout)*time.Second)
}

func CtxWithCancel(parentCtx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(parentCtx)
}

func DbFromCtx(ctx context.Context) (*db.Postgres, error) {
	d, ok := ctx.Value(dependsKey).(*depends)
	if !ok {
		return nil, ErrGetDepends
	}
	return d.pg, nil
}
