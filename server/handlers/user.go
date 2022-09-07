package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"github.com/fabienbellanger/go-url-shortener/server/db"
	models "github.com/fabienbellanger/go-url-shortener/server/models"
	"github.com/fabienbellanger/go-url-shortener/server/repositories"
	"github.com/fabienbellanger/go-url-shortener/server/utils"
	"github.com/fabienbellanger/goutils/mail"
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
		createErrors := utils.ValidateStruct(*user)
		if createErrors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
				Details: createErrors,
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

// ForgottenPassword save a forgotten password request.
func ForgottenPassword(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Find user
		user, err := repositories.GetUserByUsername(db, c.Params("email"))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when retrieving user")
		}
		if user.ID == "" {
			return c.Status(fiber.StatusNotFound).JSON(utils.HTTPError{
				Code:    fiber.StatusNotFound,
				Message: "No user found",
			})
		}

		// Sale line in database
		passwordReset := models.PasswordResets{
			UserID:    user.ID,
			Token:     uuid.New().String(),
			ExpiredAt: time.Now().Add(viper.GetDuration("FORGOTTEN_PASSWORD_EXPIRATION_DURATION") * time.Hour).UTC(),
		}
		err = repositories.CreateOrUpdatePasswordReset(db, &passwordReset)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when requesting new password")
		}

		// Send email with link
		to := make([]string, 1)
		to[0] = user.Username
		subject := fmt.Sprintf("[%s] Forgotten password", viper.GetString("APP_NAME"))
		var body bytes.Buffer

		tp, err := template.ParseFiles("templates/forgotten_password.gohtml")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when creating password reset email")
		}
		err = tp.Execute(&body, struct {
			Title string
			Link  string
		}{
			Title: fmt.Sprintf("%s - Forgotten password", viper.GetString("APP_NAME")),
			Link:  fmt.Sprintf("%s/%s", viper.GetString("FORGOTTEN_PASSWORD_BASE_URL"), passwordReset.Token),
		})
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when creating password reset email")
		}

		err = mail.Send(viper.GetString("FORGOTTEN_PASSWORD_EMAIL_FROM"), to, subject, body.String(), "", "", viper.GetString("SMTP_HOST"), viper.GetInt("SMTP_PORT"))
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when sending password reset email")
		}

		return c.JSON(passwordReset)
	}
}

// UpdateUserPassword update user password.
func UpdateUserPassword(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Params("token")

		newPassword := new(models.UserUpdatePassword)
		if err := c.BodyParser(newPassword); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		// Data validation
		// ---------------
		createErrors := utils.ValidateStruct(*newPassword)
		if createErrors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
				Details: createErrors,
			})
		}

		// Update user password
		// --------------------
		userID, err := repositories.GetUserIDFromPasswordReset(db, token, newPassword.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when searching user")
		}
		if userID == "" {
			return fiber.NewError(fiber.StatusNotFound, "no user found")
		}
		err = repositories.UpdateUserPassword(db, userID, newPassword.Password)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when updating user password")
		}

		// Delete password reset
		// ---------------------
		err = repositories.DeletePasswordReset(db, userID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when deleting user password reset")
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
