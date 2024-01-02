package handler

import (
	"errors"
	"kazokku/internal/app/delivery/dto"
	"kazokku/internal/app/service"
	"kazokku/internal/helpers"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) userHandler {
	return userHandler{userService}
}

func (h userHandler) Register(ctx *fiber.Ctx) error {
	var data dto.UserRegisterRequest
	if err := ctx.BodyParser(&data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id, err := h.userService.Create(ctx, data)
	if err != nil {
		var respErr helpers.ResponseError
		if errors.As(err, &respErr) {
			return ctx.Status(respErr.Code()).JSON(fiber.Map{
				"error": respErr.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"user_id": id,
	})
}

func (h userHandler) GetAll(ctx *fiber.Ctx) error {
	var query dto.UserQuery
	if err := ctx.QueryParser(&query); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	users, err := h.userService.GetAll(ctx, query)
	if err != nil {
		var respErr helpers.ResponseError
		if errors.As(err, &respErr) {
			return ctx.Status(respErr.Code()).JSON(fiber.Map{
				"error": respErr.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"count": len(users),
		"rows":  users,
	})
}

func (h userHandler) GetByID(ctx *fiber.Ctx) error {
	userID, err := ctx.ParamsInt("user_id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user, err := h.userService.GetByID(ctx, uint(userID))
	if err != nil {
		var respErr helpers.ResponseError
		if errors.As(err, &respErr) {
			return ctx.Status(respErr.Code()).JSON(fiber.Map{
				"error": respErr.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(user)
}

func (h userHandler) UpdateByID(ctx *fiber.Ctx) error {
	var data dto.UserUpdateRequest
	if err := ctx.BodyParser(&data); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.userService.UpdateByID(ctx, data); err != nil {
		var respErr helpers.ResponseError
		if errors.As(err, &respErr) {
			return ctx.Status(respErr.Code()).JSON(fiber.Map{
				"error": respErr.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
