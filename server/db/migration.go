package db

import "github.com/fabienbellanger/go-url-shortener/server/models"

// modelsList lists all models to automigrate.
var modelsList = []any{
	&models.User{},
	&models.Link{},
	&models.PasswordResets{},
}
