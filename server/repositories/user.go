package repositories

import (
	"crypto/sha512"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"

	"github.com/fabienbellanger/go-url-shortener/server/db"
	"github.com/fabienbellanger/go-url-shortener/server/models"
)

// Login gets user from username and password.
func Login(db *db.DB, username, password string) (user models.User, err error) {
	// Hash password
	// -------------
	passwordBytes := sha512.Sum512([]byte(password))
	password = hex.EncodeToString(passwordBytes[:])

	if result := db.Where(&models.User{Username: username, Password: password}).First(&user); result.Error != nil {
		return user, result.Error
	}
	return user, err
}

// ListAllUsers gets all users in database.
func ListAllUsers(db *db.DB) ([]models.User, error) {
	var users []models.User

	if response := db.Find(&users); response.Error != nil {
		return users, response.Error
	}
	return users, nil
}

// CreateUser adds user in database.
func CreateUser(db *db.DB, user *models.User) error {
	// UUID
	// ----
	user.ID = uuid.New().String()

	// Hash password
	// -------------
	passwordBytes := sha512.Sum512([]byte(user.Password))
	user.Password = hex.EncodeToString(passwordBytes[:])

	if result := db.Create(&user); result.Error != nil {
		return result.Error
	}
	return nil
}

// GetUser returns a user from its ID.
func GetUser(db *db.DB, id string) (user models.User, err error) {
	if result := db.Find(&user, "id = ?", id); result.Error != nil {
		return user, result.Error
	}
	return user, err
}

// GetUser returns a user from its username.
func GetUserByUsername(db *db.DB, username string) (user models.User, err error) {
	if result := db.Find(&user, "username = ?", username); result.Error != nil {
		return user, result.Error
	}
	return user, err
}

// DeleteUser deletes a user from database.
func DeleteUser(db *db.DB, id string) error {
	result := db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UpdateUser updates user information.
func UpdateUser(db *db.DB, id string, userForm *models.UserForm) (user models.User, err error) {
	// Hash password
	// -------------
	hashedPassword := sha512.Sum512([]byte(userForm.Password))

	result := db.Model(&models.User{}).Where("id = ?", id).Select("lastname", "firstname", "username", "password").Updates(models.User{
		Lastname:  userForm.Lastname,
		Firstname: userForm.Firstname,
		Username:  userForm.Username,
		Password:  hex.EncodeToString(hashedPassword[:]),
	})
	if result.Error != nil {
		return user, result.Error
	}

	user, err = GetUser(db, id)
	if err != nil {
		return user, err
	}
	return user, err
}

// UpdateUserPassword updates user passwords.
func UpdateUserPassword(db *db.DB, id, password string) error {
	// Hash password
	// -------------
	hashedPassword := sha512.Sum512([]byte(password))

	result := db.Exec(`
		UPDATE users
		SET password = ?, updated_at = ?
		WHERE id = ?`,
		hex.EncodeToString(hashedPassword[:]),
		time.Now().UTC(),
		id,
	)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// CreatePasswordReset add a reset password request in database.
func CreatePasswordReset(db *db.DB, passwordReset *models.PasswordResets) error {
	result := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&passwordReset)

	return result.Error
}

// GetUserIDFromPasswordReset update user password and delete password_resets line.
func GetUserIDFromPasswordReset(db *db.DB, token, password string) (string, error) {
	data := struct {
		ID string
	}{}

	result := db.Raw(`
			SELECT u.id AS id
			FROM password_resets pr
				INNER JOIN users u ON u.id = pr.user_id AND u.deleted_at IS NULL
			WHERE pr.token = ?
				AND pr.expired_at >= ?`,
		token,
		time.Now().UTC()).Scan(&data)
	if result.Error != nil {
		return "", result.Error
	}

	return data.ID, nil
}

// DeletePasswordReset deletes user password reset.
func DeletePasswordReset(db *db.DB, userId string) error {
	result := db.Where("user_id = ?", userId).Delete(&models.PasswordResets{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
