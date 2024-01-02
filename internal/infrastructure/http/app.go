package http

import (
	"fmt"
	"kazokku/internal/app/delivery/routes"
	"kazokku/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type App struct {
	app  *fiber.App
	host string
	port int
}

func New(conf utils.App, db *pgx.Conn) App {
	app := fiber.New()
	routes.NewUserRoutes(db, app)

	app.Static("/photos", conf.SaveDir)

	return App{
		app:  app,
		host: conf.Host,
		port: conf.Port,
	}
}

func (a App) Run() error {
	return a.app.Listen(fmt.Sprintf("%s:%d", a.host, a.port))
}
