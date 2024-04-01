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
	dbString := "postgres://" + cfg.Database.Username + ":" + cfg.Database.Password + "@" + cfg.Database.Host + ":" + cfg.Database.Port + "/" + cfg.Database.DBName + "?sslmode=disable"
	dbConn, err := db.NewPG(context.Background(), dbString)
	if err != nil {
		log.Printf("ERROR | failed to connect to database with error: %s\n", err)
	}
	ctx := tools.CtxWithDepends(context.Background(), dbConn)
	repo := new(repositories.DbRepository)
	services.AllItems(ctx, cfg.Chrome, repo)
}
