package server

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/fabienbellanger/go-url-shortener/db"
	"github.com/fabienbellanger/go-url-shortener/handlers"
)

// Web routes
// ----------

func registerPublicWebRoutes(r fiber.Router, db *db.DB, logger *zap.Logger) {
	// Admin interface
	registerPublicAdminRoutes(r, db, logger)

	// Shorted URL
	r.Get("/:id", handlers.RedirectURL(db, logger))
}

func registerPublicAdminRoutes(r fiber.Router, db *db.DB, logger *zap.Logger) {
	// TODO:
	//   - https://github.com/gofiber/template/tree/master/html
	//   - https://docs.gofiber.io/guide/templates
	//   - https://docs.gofiber.io/api/middleware/session#examples
	admin := r.Group("admin")

	admin.Get("/", func(c *fiber.Ctx) error {
		return c.Render("public/admin/index", fiber.Map{
			"Title": "Admin interface",
		})
	})

	admin.Get("/login", handlers.GetLoginPage()).Name("loginPage")
	admin.Post("/login", handlers.PostLoginPage())
}

// API routes
// ----------

func registerPublicAPIRoutes(r fiber.Router, db *db.DB) {
	v1 := r.Group("/v1")

	v1.Post("/login", handlers.Login(db))
}

func registerProtectedAPIRoutes(r fiber.Router, db *db.DB) {
	v1 := r.Group("/v1")

	// Register
	v1.Post("/register", handlers.CreateUser(db))

	// Users
	users := v1.Group("/users")
	users.Get("/", handlers.GetAllUsers(db))
	users.Get("/:id", handlers.GetUser(db))
	users.Delete("/:id", handlers.DeleteUser(db))
	users.Put("/:id", handlers.UpdateUser(db))

	// Links
	links := v1.Group("/links")
	links.Get("", handlers.LinksList(db))
	links.Post("", handlers.CreateLink(db))
}
