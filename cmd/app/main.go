package main

import (
	"context"
	"kazokku/internal/infrastructure/database"
	"kazokku/internal/infrastructure/http"
	"kazokku/internal/utils"
)

func main() {
	conf, err := utils.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	err = database.Migrate(conf.Database)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	db, err := database.New(ctx, conf.Database)
	if err != nil {
		panic(err)
	}
	defer db.Close(ctx)

	app := http.New(conf.App, db)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
