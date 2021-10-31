package handlers

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/fabienbellanger/go-url-shortener/db"
	models "github.com/fabienbellanger/go-url-shortener/models"
	"github.com/fabienbellanger/go-url-shortener/repositories"
	"github.com/fabienbellanger/go-url-shortener/utils"
)

type userLogin struct {
	models.User
	Token     string `json:"token" xml:"token" form:"token"`
	ExpiresAt string `json:"expires_at" xml:"expires_at" form:"expires_at"`
}

type userAuth struct {
	Username string `json:"username" xml:"username" form:"username" validate:"required,email"`
	Password string `json:"password" xml:"password" form:"password" validate:"required,min=8"`
}

// Login authenticates a user.
func Login(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		u := new(userAuth)
		if err := c.BodyParser(u); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		loginErrors := utils.ValidateStruct(*u)
		if loginErrors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
				Details: loginErrors,
			})
		}

		user, err := repositories.Login(db, u.Username, u.Password)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusUnauthorized).JSON(utils.HTTPError{
					Code:    fiber.StatusUnauthorized,
					Message: "Unauthorized",
				})
			}
			return fiber.NewError(fiber.StatusInternalServerError, "Error during authentication")
		}

		// Create token
		token := jwt.New(jwt.SigningMethodHS512)

		// Expiration time
		now := time.Now()
		expiresAt := now.Add(time.Hour * viper.GetDuration("JWT_LIFETIME"))

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = user.ID
		claims["username"] = user.Username
		claims["lastname"] = user.Lastname
		claims["firstname"] = user.Firstname
		claims["createdAt"] = user.CreatedAt
		claims["exp"] = expiresAt.Unix()
		claims["iat"] = now.Unix()
		claims["nbf"] = now.Unix()

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(viper.GetString("JWT_SECRET")))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error during token generation")
		}

		return c.JSON(userLogin{
			User:      user,
			Token:     t,
			ExpiresAt: expiresAt.Format("2006-01-02T15:04:05.000Z"),
		})
	}
}

// GetAllUsers lists all users.
func GetAllUsers(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		users, err := repositories.ListAllUsers(db)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(users)
	}
}

// GetUser return a user.
func GetUser(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad ID",
			})
		}

		user, err := repositories.GetUser(db, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when retrieving user")
		}
		if user.ID == "" {
			return c.Status(fiber.StatusNotFound).JSON(utils.HTTPError{
				Code:    fiber.StatusNotFound,
				Message: "No user found",
			})
		}

		return c.JSON(user)
	}
}

// CreateUser creates a new user.
func CreateUser(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(models.UserForm)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		// Data validation
		// ---------------
		if user.Firstname == "" || user.Lastname == "" || user.Username == "" || user.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Parameters",
			})
		}

		// Database insertion
		// ------------------
		newUser := models.User{
			Lastname:  user.Lastname,
			Firstname: user.Firstname,
			Password:  user.Password,
			Username:  user.Username,
		}

		if err := repositories.CreateUser(db, &newUser); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error during user creation")
		}
		return c.JSON(newUser)
	}
}

// DeleteUser return a user.
func DeleteUser(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad ID",
			})
		}

		err := repositories.DeleteUser(db, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when deleting user")
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// UpdateUser updates user information.
func UpdateUser(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad ID",
			})
		}

		user := new(models.UserForm)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Data",
			})
		}

		updateErrors := utils.ValidateStruct(*user)
		if updateErrors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
				Details: updateErrors,
			})
		}

		updatedUser, err := repositories.UpdateUser(db, id, user)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when updating user")
		}

		return c.JSON(updatedUser)
	}
}
