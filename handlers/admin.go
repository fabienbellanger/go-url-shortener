package handlers

import (
	"github.com/fabienbellanger/go-url-shortener/db"
	"github.com/fabienbellanger/go-url-shortener/repositories"
	"github.com/fabienbellanger/go-url-shortener/utils"
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
//
// TODO: Send errors to login page
func PostLoginPage(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		u := new(userAuth)
		if err := c.BodyParser(u); err != nil {
			return c.RedirectToRoute("loginPage", nil)
		}

		loginErrors := utils.ValidateStruct(*u)
		if loginErrors != nil {
			return c.RedirectToRoute("loginPage", nil)
		}

		_, err := repositories.Login(db, u.Username, u.Password)
		if err != nil {
			return c.RedirectToRoute("loginPage", nil)
		}

		return c.RedirectToRoute("linksPage", nil)
	}
}
