package db

import "github.com/fabienbellanger/go-url-shortener/models"

// modelsList lists all models to automigrate.
var modelsList = []interface{}{
	&models.User{},
	&models.Link{},
}
