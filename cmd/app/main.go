package main

import (
	"github.com/oustrix/ozon_journal/config"
	"github.com/oustrix/ozon_journal/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	app.Run(cfg)
}
