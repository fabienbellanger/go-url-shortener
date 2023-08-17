package server

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/fabienbellanger/go-url-shortener/server/db"
	"github.com/fabienbellanger/go-url-shortener/server/handlers"
)

// Web routes
// ----------

func registerPublicWebRoutes(r fiber.Router, db *db.DB, logger *zap.Logger) {
	// Basic Auth
	// ----------
	cfg := basicauth.Config{
		Users: map[string]string{
			viper.GetString("SERVER_BASICAUTH_USERNAME"): viper.GetString("SERVER_BASICAUTH_PASSWORD"),
		},
	}

	// API documentation
	doc := r.Group("/doc")
	doc.Use(basicauth.New(cfg))
	doc.Get("/api-v1", handlers.DocAPIv1())

	// Shorted URL
	r.Get("/:id", handlers.RedirectURL(db, logger))

	// Filesystem
	// ----------
	assets := r.Group("/assets")
	assets.Use(filesystem.New(filesystem.Config{
		Root:   http.Dir("./assets"),
		Browse: false,
		Index:  "index.html",
		MaxAge: 3600,
	}))
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

func registerProtectedAPIRoutes(r fiber.Router, db *db.DB, logger *zap.Logger) {
	v1 := r.Group("/v1")

	// Register
	v1.Post("/register", handlers.CreateUser(db))

	// Users
	users := v1.Group("/users")
	users.Get("", handlers.GetAllUsers(db))
	users.Get("/:id", handlers.GetUser(db))
	users.Put("/:id", handlers.UpdateUser(db))
	users.Delete("/:id", handlers.DeleteUser(db))

	// Links
	links := v1.Group("/links")
	links.Get("", handlers.LinksList(db))
	links.Post("", handlers.CreateLink(db))
	links.Post("/upload", handlers.UploadLink(db, logger))
	links.Put("/:id", handlers.UpdateLink(db))
	links.Delete("/selected", handlers.DeleteLinks(db))
	links.Delete("/:id", handlers.DeleteLink(db))
	links.Get("/export/csv", handlers.ExportCSVLinks(db, logger))
}
