package main

import (
	"context"
	"github.com/L1LSunflower/lina/config"
	"github.com/L1LSunflower/lina/internal/repositories"
	"github.com/L1LSunflower/lina/internal/services"
	"github.com/L1LSunflower/lina/internal/tools"
	"github.com/L1LSunflower/lina/pkg/db"
	"log"
)

func main() {
	cfg := config.Conf()
	dbConn, err := db.NewPG(context.Background(), cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode)
	if err != nil {
		log.Printf("ERROR | failed to connect to database with error: %s\n", err)
	}
	ctx := tools.CtxWithDepends(context.Background(), dbConn)
	repo := new(repositories.ItemsRepository)
	services.AllItems(ctx, cfg.Chrome, repo)
}
