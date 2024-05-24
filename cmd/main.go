package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/L1LSunflower/lina/config"

	"github.com/L1LSunflower/lina/internal/services"
	"github.com/L1LSunflower/lina/internal/tools"

	"github.com/L1LSunflower/lina/pkg/db"
	"github.com/L1LSunflower/lina/pkg/telegram"
)

func main() {
	cfg := config.Conf()
	dbConn, err := db.NewPG(context.Background(), cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)
	if err != nil {
		log.Printf("ERROR | failed to connect to database with error: %s\n", err)
	}
	queue := tools.NewQueue(cfg.Workers)
	go queue.Run()
	// Get all items
	itemCollector := services.NewItemCollector(cfg, dbConn, queue)
	go itemCollector.Start()
	// Telegram handler
	tHandler := telegram.NewTHandler(cfg.Telegram.TelegramUrl, cfg.Telegram.TelegramKey)
	// Workers
	workers := services.NewWorker(cfg, dbConn, tHandler, queue)
	// Start workers
	go workers.Start()
	log.Println("application start")
	// Wait workers
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	// Stopping workers
	log.Println("shutting down...")
	itemCollector.Stop()
	workers.Stop()
	queue.Stop()
	dbConn.Close()
}
