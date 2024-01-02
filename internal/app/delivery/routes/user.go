package routes

import (
	"kazokku/internal/app/delivery/handler"
	"kazokku/internal/app/delivery/middleware"
	"kazokku/internal/app/repository"
	"kazokku/internal/app/service"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func NewUserRoutes(db *pgx.Conn, app *fiber.App) {
	userRepo := repository.NewUserRepository(db)
	ccRepo := repository.NewCreditCardRepository(db)
	photoRepo := repository.NewPhotoRepository(db)
	userService := service.NewUserService(db, userRepo, ccRepo, photoRepo)
	userHandler := handler.NewUserHandler(userService)
	user := app.Group("/user")

	user.Use(middleware.ApiKey())
	{
		user.Post("/register", userHandler.Register)
		user.Get("/list", userHandler.GetAll)
		user.Get("/:user_id", userHandler.GetByID)
		user.Patch("", userHandler.UpdateByID)
	}
}
