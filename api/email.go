package api

import "github.com/gofiber/fiber/v3"

func CreateNewEmailHandler(c fiber.Ctx) error {
	return c.JSON("{ \"name\": true }")
}
