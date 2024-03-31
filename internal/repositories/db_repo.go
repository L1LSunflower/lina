package repositories

import (
	"context"

	"github.com/L1LSunflower/lina/internal/entities"
)

type DBRepo interface {
	AddItem(ctx context.Context, item *entities.Item) error
	AddItems(ctx context.Context, items []*entities.Item) error
	Items(ctx context.Context, id, status string, limit int) ([]*entities.Item, error)
	CheckByHash(ctx context.Context, hash string) (bool, error)
}
