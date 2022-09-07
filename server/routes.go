package server

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/fabienbellanger/go-url-shortener/server/db"
	"github.com/fabienbellanger/go-url-shortener/server/handlers"
)

// Web routes
// ----------

func registerPublicWebRoutes(r fiber.Router, db *db.DB, logger *zap.Logger) {
	// Shorted URL
	r.Get("/:id", handlers.RedirectURL(db, logger))
}

// API routes
// ----------

func registerPublicAPIRoutes(r fiber.Router, db *db.DB) {
	v1 := r.Group("/v1")

	v1.Post("/login", handlers.Login(db))

	// Password reset
	v1.Post("/forgotten-password/:email", handlers.ForgottenPassword(db))
	v1.Patch("/update-password/:token", handlers.UpdateUserPassword(db))
}

func registerProtectedAPIRoutes(r fiber.Router, db *db.DB) {
	v1 := r.Group("/v1")

	// Register
	v1.Post("/register", handlers.CreateUser(db))

	// Users
	users := v1.Group("/users")
	users.Get("/", handlers.GetAllUsers(db))
	users.Get("/:id", handlers.GetUser(db))
	users.Put("/:id", handlers.UpdateUser(db))
	users.Delete("/:id", handlers.DeleteUser(db))

	// Links
	links := v1.Group("/links")
	links.Get("", handlers.LinksList(db))
	links.Post("", handlers.CreateLink(db))
	links.Put("/:id", handlers.UpdateLink(db))
	links.Delete("/:id", handlers.DeleteLink(db))
}
