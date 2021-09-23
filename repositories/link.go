package repositories

import (
	"errors"
	"time"

	"github.com/fabienbellanger/go-url-shortener/db"
	"github.com/fabienbellanger/go-url-shortener/models"
	"github.com/fabienbellanger/go-url-shortener/utils"
	"github.com/spf13/viper"
)

func getLinksFromURL(db *db.DB, url string) (links []models.Link, err error) {
	if result := db.Find(&links, "url = ?", url); result.Error != nil {
		return links, result.Error
	}
	return
}

// CreateLink adds a shortened URL in database.
func CreateLink(db *db.DB, link *models.LinkForm) (newLink models.Link, err error) {
	// Check if original URL is not already in database.
	links, err := getLinksFromURL(db, link.URL)
	if err != nil {
		return newLink, errors.New("too many links with the same URL")
	} else if len(links) > 1 {
		return newLink, errors.New("too many links with the same URL")
	} else if len(links) == 1 {
		updatedLink := links[0]
		if link.ExpiredAt != updatedLink.ExpiredAt {
			// Update expired datetime
			// -----------------------
			result := db.Model(&updatedLink).Update("expired_at", link.ExpiredAt)
			if result.Error != nil {
				return newLink, result.Error
			}
		}
		return updatedLink, nil
	}

	// Create new link
	// ---------------
	newLink = models.Link{
		ID:        utils.GenerateShortLink(link.URL, viper.GetString("URL_ENCODER_KEY")),
		URL:       link.URL,
		ExpiredAt: link.ExpiredAt,
	}

	if result := db.Create(&newLink); result.Error != nil {
		return newLink, result.Error
	}
	return
}

// GetLinkFromID returns a link if ID exists, else returns an error.
func GetLinkFromID(db *db.DB, id string) (link *models.Link, err error) {
	result := db.Where("expired_at >= ?", time.Now()).First(&link, "id = ?", id)
	if result.Error != nil {
		return link, result.Error
	}
	return
}
