package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/L1LSunflower/lina/internal/entities"
	"github.com/L1LSunflower/lina/internal/tools"
)

const (
	delimiter      = ","
	defaultTimeout = 5
)

type ItemsRepository struct{}

func NewItemsRepository() Items {
	return &ItemsRepository{}
}

func (d *ItemsRepository) AddItem(ctx context.Context, item *entities.Item) error {
	db, err := tools.DbFromCtx(ctx)
	if err != nil {
		return err
	}
	itemArgs := pgx.NamedArgs{
		"id":            item.ID,
		"url":           item.URL,
		"name":          item.Name,
		"article":       item.Article,
		"price":         item.ExpectedPrice,
		"price_on_sale": item.ActualPrice,
		"currency":      item.Currency,
		"colors":        item.Colors,
		"sizes":         item.Sizes,
		"image_links":   item.ImageLinks,
		"hash":          item.Hash,
		"status":        item.Status,
	}
	fieldsInLine, namedArgs := FieldsAndArgsFromMap(itemArgs, delimiter)
	dbCtx, cancelFunc := tools.CtxWithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()
	query := "insert into items (" + fieldsInLine + ") values (" + namedArgs + ")"
	if _, err = db.Pool.Exec(dbCtx, query, itemArgs); err != nil {
		return fmt.Errorf("unable to insert row: %w", err)
	}
	return nil
}

func (d *ItemsRepository) AddItems(ctx context.Context, items []*entities.Item) error {
	db, err := tools.DbFromCtx(ctx)
	if err != nil {
		return err
	}
	fields := []string{
		"id",
		"url",
		"name",
		"article",
		"price",
		"price_on_sale",
		"currency",
		"colors",
		"sizes",
		"image_links",
		"hash",
		"status",
	}
	fieldsInLine, namedArgs := FieldsAndArgsFromSlice(fields, delimiter)
	dbCtx, cancelFunc := tools.CtxWithTimeout(context.Background(), 20 /*defaultTimeout*/)
	defer cancelFunc()
	query := "insert into items (" + fieldsInLine + ") values (" + namedArgs + ")"
	batch := new(pgx.Batch)
	for _, item := range items {
		args := pgx.NamedArgs{
			"id":            item.ID,
			"url":           item.URL,
			"name":          item.Name,
			"article":       item.Article,
			"price":         item.ExpectedPrice,
			"price_on_sale": item.ActualPrice,
			"currency":      item.Currency,
			"colors":        item.Colors,
			"sizes":         item.Sizes,
			"image_links":   item.ImageLinks,
			"hash":          item.Hash,
			"status":        item.Status,
		}
		batch.Queue(query, args)
	}
	results := db.Pool.SendBatch(dbCtx, batch)
	defer results.Close()
	for _, item := range items {
		if _, err = results.Exec(); err != nil {
			fmt.Println(fmt.Errorf("unable to insert row by name: %s with error: %w", item.Name, err))
		}
	}
	return nil
}

func (d *ItemsRepository) Items(ctx context.Context, id, status string, limit int) ([]*entities.Item, error) {
	db, err := tools.DbFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	fields := []string{
		"id",
		"url",
		"name",
		"article",
		"price",
		"price_on_sale",
		"currency",
		"colors",
		"sizes",
		"image_links",
		"hash",
		"status",
	}
	query := "select " + Fields(fields, delimiter) + " from items where id<@id and status=@status order by id limit @limit"
	args := pgx.NamedArgs{
		"id":     id,
		"limit":  limit,
		"status": status,
	}
	dbCtx, cancelFunc := tools.CtxWithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()
	rows, err := db.Pool.Query(dbCtx, query, args)
	if err != nil {
		return nil, err
	}
	var items []*entities.Item
	for rows.Next() {
		item := new(entities.Item)
		if err = rows.Scan(
			&item.ID,
			&item.URL,
			&item.Name,
			&item.Article,
			&item.ExpectedPrice,
			&item.ActualPrice,
			&item.Currency,
			&item.Colors,
			&item.Sizes,
			&item.ImageLinks,
			&item.Hash,
			&item.Status,
		); err != nil {
			return nil, fmt.Errorf("unable to scan row %w", err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (d *ItemsRepository) CheckByHash(ctx context.Context, hash string) (bool, error) {
	db, err := tools.DbFromCtx(ctx)
	if err != nil {
		return false, err
	}
	query := "select exists(select id from items where hash=@hash)"
	args := pgx.NamedArgs{"hash": hash}
	dbCtx, cancelFunc := tools.CtxWithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()
	exist := false
	if err = db.Pool.QueryRow(dbCtx, query, args).Scan(&exist); err != nil {
		return false, err
	}
	return exist, nil
}

func (d *ItemsRepository) UpdateStatus(ctx context.Context, id, status string) error {
	db, err := tools.DbFromCtx(ctx)
	if err != nil {
		return err
	}
	query := "update items set status=@status where id=@id"
	dbCtx, cancelFunc := tools.CtxWithTimeout(context.Background(), defaultTimeout)
	defer cancelFunc()
	if _, err = db.Pool.Exec(dbCtx, query, pgx.NamedArgs{"id": id, "status": status}); err != nil {
		return err
	}
	return nil
}
