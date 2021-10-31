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
	r.Get("/:id", handlers.RedirectURL(db))
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
