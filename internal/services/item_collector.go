package services

import (
	"context"
	"encoding/json"
	"github.com/L1LSunflower/lina/internal/drivers/lichi"
	"log"
	"time"

	"github.com/L1LSunflower/lina/config"

	"github.com/L1LSunflower/lina/internal/entities"
	"github.com/L1LSunflower/lina/internal/repositories"

	"github.com/L1LSunflower/lina/internal/tools"
	"github.com/L1LSunflower/lina/pkg/db"
	"github.com/L1LSunflower/lina/pkg/encryption"
)

const (
	defaultTimeoutForSaver = 30 * time.Second
	defaultTimeoutForLichi = 1 * time.Hour
)

type ItemCollector struct {
	cfg             *config.Chrome
	dbConn          *db.Postgres
	queue           tools.Jobs
	itemsChan       chan *entities.Item
	itemsToSaveChan chan *entities.Item
	cancelChan      chan struct{}
	repo            repositories.Items
}

func NewItemCollector(cfg *config.Config, dbConn *db.Postgres, queue tools.Jobs) *ItemCollector {
	return &ItemCollector{
		cfg:             &cfg.Chrome,
		dbConn:          dbConn,
		queue:           queue,
		itemsChan:       make(chan *entities.Item, cfg.Workers),
		itemsToSaveChan: make(chan *entities.Item, cfg.Workers),
		cancelChan:      make(chan struct{}, 1),
		repo:            repositories.NewItemsRepository(),
	}
}

func (i *ItemCollector) Start() {
	jobTicker := time.NewTicker(defaultTimeoutForSaver)
	dataTicker := time.NewTicker(defaultTimeoutForLichi)
	dataSaveTicker := time.NewTicker(defaultTimeoutForSaver)
	for {
		select {
		case <-jobTicker.C:
			ctx := tools.CtxWithDepends(context.Background(), i.dbConn)
			i.queue.GoWithTimeoutWithContext(ctx, defaultTimeoutForJob, i.processItem)
		case <-dataTicker.C:
			ctx := tools.CtxWithDepends(context.Background(), i.dbConn)
			i.queue.GoWithTimeoutWithContext(ctx, defaultTimeoutForExporter, i.parseItems)
		case <-dataSaveTicker.C:
			ctx := tools.CtxWithDepends(context.Background(), i.dbConn)
			i.queue.GoWithTimeoutWithContext(ctx, defaultTimeoutForExporter, i.saveItem)
		case <-i.cancelChan:
			i.Stop()
			break
		}
	}
}

func (i *ItemCollector) Stop() {
	close(i.cancelChan)
	close(i.itemsChan)
	close(i.itemsToSaveChan)
}

func (i *ItemCollector) saveItem(ctx context.Context) {
	log.Println("SAVING ITEM START")
	select {
	case item := <-i.itemsToSaveChan:
		if err := i.repo.AddItem(ctx, item); err != nil {
			log.Printf("ERROR | failed to save item with error: %s\n", err)
		}
		log.Println("SAVING ITEM END")
	case <-ctx.Done():
		return
	}
}

func (i *ItemCollector) parseItems(ctx context.Context) {
	log.Println("PARSING ITEMS START")
	select {
	case <-ctx.Done():
		return
	default:
		ctx, cancel := tools.CtxWithCancel(ctx)
		defer cancel()
		items, err := lichi.Items(ctx, i.cfg.Headless, i.cfg.DebugMode)
		if err != nil {
			log.Println(" | ERROR: Get new item with error:", err)
		}
		if len(items) == 0 {
			return
		}
		if _, isChannelOpen := <-i.cancelChan; !isChannelOpen {
			return
		}
		for _, item := range items {
			if _, isChannelOpen := <-i.cancelChan; !isChannelOpen {
				return
			}
			i.itemsChan <- item
		}
		log.Println("PARSING ITEMS END")
	}
}

func (i *ItemCollector) processItem(ctx context.Context) {
	log.Println("PROCESSING ITEM")
	select {
	case item := <-i.itemsChan:
		sanitizedItem := item.SanitizedItem()
		itemB, err := json.Marshal(sanitizedItem)
		if err != nil {
			log.Printf("ERROR | error while marshaling item by url: %s | with error: %s\n", item.URL, err)
		}
		itemHash := encryption.Hash(itemB)
		exist, err := i.repo.CheckByHash(ctx, itemHash)
		if err != nil {
			log.Printf("ERROR | failed to check item by hash with url: %s | with error: %s\n", item.URL, err)
		}
		if exist {
			return
		}
		item.PrepareToSave(itemHash)
		if _, isChannelOpen := <-i.cancelChan; !isChannelOpen {
			return
		}
		i.itemsToSaveChan <- item
		log.Printf("item: %v\n", item)
		log.Println("PROCESSING END")
	case <-ctx.Done():
		return
	}
}
