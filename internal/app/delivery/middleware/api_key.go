package middleware

import "github.com/gofiber/fiber/v2"

func ApiKey() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if apiKey := string(c.Request().Header.Peek("key")); apiKey == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "API Key is missing.",
			})
		} else if apiKey != "HiJhvL$T27@1u^%u86g" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API Key.",
			})
		}
		return c.Next()
	}
}
