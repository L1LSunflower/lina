package main

import (
	"github.com/L1LSunflower/lina/config"
	"github.com/L1LSunflower/lina/internal/services"
)

func main() {
	cfg := config.Conf()
	services.AllItems(cfg)
}
