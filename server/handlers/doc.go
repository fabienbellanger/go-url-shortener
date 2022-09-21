package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// DocAPIv1 show API v1 documentation.
func DocAPIv1() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("doc_api_v1", fiber.Map{})
	}
}
