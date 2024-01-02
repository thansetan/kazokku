package http

import (
	"fmt"
	"kazokku/internal/app/delivery/routes"
	"kazokku/internal/utils"
	"log/slog"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	loggerMW "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	app  *fiber.App
	host string
	port int
}

func New(conf utils.App, db *pgxpool.Pool, logger *slog.Logger) App {
	app := fiber.New()
	app.Use(recover.New())
	app.Use(loggerMW.New())
	app.Use(requestid.New())

	routes.NewUserRoutes(db, app, logger)

	app.Static("/photos", filepath.Join(conf.SaveDir, "photos"))

	return App{
		app:  app,
		host: conf.Host,
		port: conf.Port,
	}
}

func (a App) Run() error {
	return a.app.Listen(fmt.Sprintf("%s:%d", a.host, a.port))
}
