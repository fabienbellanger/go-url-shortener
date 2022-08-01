package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// GetLoginPage serve login page.
func GetLoginPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Render("public/admin/login", fiber.Map{
			"Title":     "Authentification",
			"csrfToken": c.Locals("csrf_token"),
		})
	}
}

// PostLoginPage
func PostLoginPage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Printf("%v\n", string(c.Body()))
		return c.RedirectToRoute("loginPage", nil)
	}
}
