package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/L1LSunflower/lina/config"

	"github.com/L1LSunflower/lina/internal/drivers/lichi"
	"github.com/L1LSunflower/lina/internal/entities"
	"github.com/L1LSunflower/lina/internal/repositories"

	"github.com/L1LSunflower/lina/internal/tools"
	"github.com/L1LSunflower/lina/pkg/encryption"
)

func AllItems(ctx context.Context, ccfg config.Chrome, repo repositories.DBRepo) {
	ctx, cancel := tools.CtxWithCancel(ctx)
	defer cancel()
	// TODO: need item: image links
	items, err := lichi.Items(ctx, ccfg.Headless, ccfg.DebugMode)
	if err != nil {
		log.Println(" | ERROR: Get new item with error:", err)
	}
	log.Println("PARSED ITEMS:")
	var (
		itemB       []byte
		exist       bool
		itemsToSave = make([]*entities.Item, 0, len(items))
	)
	for id, item := range items {
		sanitizedItem := item.SanitizedItem()
		itemB, err = json.Marshal(sanitizedItem)
		if err != nil {
			log.Printf("ERROR | error while marshaling item by url: %s | with error: %s\n", item.Url, err)
		}
		itemHash := encryption.Hash(itemB)
		if exist, err = repo.CheckByHash(ctx, itemHash); err != nil {
			log.Printf("ERROR | failed to check item by hash with url: %s | with error: %s\n", item.Url, err)
		}
		if !exist {
			item.PrepareToSave(itemHash)
			itemsToSave = append(itemsToSave, item)
		}
		log.Printf("id: %d | item: %v\n", id, item)
	}
	if len(itemsToSave) <= 0 {
		log.Printf("WARN | no items to save")
		return
	}
	if err = repo.AddItems(ctx, itemsToSave); err != nil {
		log.Printf("ERROR | failed to save item with error: %s\n", err)
	}
}
