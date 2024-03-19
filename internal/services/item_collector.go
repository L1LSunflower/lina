package services

import (
	"context"
	"log"

	"github.com/L1LSunflower/lina/config"
	"github.com/L1LSunflower/lina/internal/drivers/lichi"
)

func AllItems(cfg *config.Config) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	items, err := lichi.Items(ctx, cfg.Chrome.Headless, cfg.Chrome.DebugMode)
	if err != nil {
		log.Println(" | ERROR: Get new item with error:", err)
	}
	log.Println("PARSED ITEMS:")
	for id, item := range items {
		log.Printf("id: %d | item: %v\n", id, item)
	}
}
