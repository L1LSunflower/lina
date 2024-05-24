package repositories

import (
	"context"

	"github.com/L1LSunflower/lina/internal/entities"
)

type Items interface {
	AddItem(ctx context.Context, item *entities.Item) error
	AddItems(ctx context.Context, items []*entities.Item) error
	Items(ctx context.Context, id, status string, limit int) ([]*entities.Item, error)
	CheckByHash(ctx context.Context, hash string) (bool, error)
	UpdateStatus(ctx context.Context, id, status string) error
}

type Users interface {
	Add(ctx context.Context, user *entities.User) error
	GetAll(ctx context.Context) ([]*entities.User, error)
	Delete(ctx context.Context, id string) error
	Block(ctx context.Context, id string) error
}
