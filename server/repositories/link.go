package repositories

import (
	"errors"
	"time"

	database "github.com/fabienbellanger/go-url-shortener/server/db"
	"github.com/fabienbellanger/go-url-shortener/server/models"
	"github.com/fabienbellanger/go-url-shortener/server/utils"
	"github.com/spf13/viper"
)

func getLinksFromURL(db *database.DB, url string) (links []models.Link, err error) {
	if result := db.Find(&links, "url = ?", url); result.Error != nil {
		return links, result.Error
	}
	return
}

func GetAllLinks(db *database.DB, page, limit string) (links []models.Link, err error) {
	var total int64
	db.Model(&links).Count(&total)

	if result := db.Scopes(database.Paginate(page, limit)).Find(&links); result.Error != nil {
		return links, result.Error
	}
	return
}

// CreateLink adds a shortened URL in database.
func CreateLink(db *database.DB, link *models.LinkForm) (newLink models.Link, err error) {
	// Check if original URL is not already in database.
	links, err := getLinksFromURL(db, link.URL)

	if err != nil {
		return newLink, errors.New("error when searching links with this URL")
	} else if len(links) > 1 {
		return newLink, errors.New("too many links with the same URL")
	} else if len(links) == 1 {
		return newLink, errors.New("a link with the same URL already exists")
	}

	// Create new link
	// ---------------
	id, err := utils.GenerateShortLink(link.URL, viper.GetString("URL_ENCODER_KEY"))
	if err != nil {
		return models.Link{}, err
	}
	newLink = models.Link{
		ID:        id,
		URL:       link.URL,
		ExpiredAt: link.ExpiredAt,
	}

	if result := db.Create(&newLink); result.Error != nil {
		return newLink, result.Error
	}
	return
}

func UpdateLink(db *database.DB, link *models.Link) error {
	r := db.Save(&link)
	if r.Error != nil {
		return r.Error
	}
	return nil
}

// GetLinkFromID returns a link if ID exists, else returns an error.
func GetLinkFromID(db *database.DB, id string) (link *models.Link, err error) {
	result := db.Where("expired_at >= ?", time.Now()).First(&link, "id = ?", id)
	if result.Error != nil {
		return link, result.Error
	}
	return
}

// DeleteLink remove the link in database.
func DeleteLink(db *database.DB, id string) error {
	r := db.Where("id = ?", id).Delete(&models.Link{})
	if r.Error != nil {
		return r.Error
	}
	return nil
}
