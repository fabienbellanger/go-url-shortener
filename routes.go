package server

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/fabienbellanger/go-url-shortener/db"
	"github.com/fabienbellanger/go-url-shortener/handlers/api"
	"github.com/fabienbellanger/go-url-shortener/handlers/web"
)

// Web routes
// ----------

func registerPublicWebRoutes(r fiber.Router, logger *zap.Logger) {
	r.Get("/health-check", web.HealthCheck(logger))
	r.Get("/hello/:name", web.Hello())
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
}

func registerAuth(r fiber.Router, db *db.DB) {
	r.Post("/login", api.Login(db))
	r.Post("/register", api.CreateUser(db))
}

func registerUser(r fiber.Router, db *db.DB) {
	users := r.Group("/users")

	users.Get("/", api.GetAllUsers(db))
	users.Get("/stream", api.StreamUsers(db))
	users.Get("/:id", api.GetUser(db))
	users.Delete("/:id", api.DeleteUser(db))
	users.Put("/:id", api.UpdateUser(db))
}
