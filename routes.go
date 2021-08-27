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

	registerAuth(v1, db)
}

func registerProtectedAPIRoutes(r fiber.Router, db *db.DB) {
	v1 := r.Group("/v1")

	registerUser(v1, db)

	// Links
	v1.Post("/links", handlers.CreateLink(db))
}

func registerAuth(r fiber.Router, db *db.DB) {
	r.Post("/login", handlers.Login(db))
	r.Post("/register", handlers.CreateUser(db))
}

func registerUser(r fiber.Router, db *db.DB) {
	users := r.Group("/users")

	users.Get("/", handlers.GetAllUsers(db))
	users.Get("/stream", handlers.StreamUsers(db))
	users.Get("/:id", handlers.GetUser(db))
	users.Delete("/:id", handlers.DeleteUser(db))
	users.Put("/:id", handlers.UpdateUser(db))
}
