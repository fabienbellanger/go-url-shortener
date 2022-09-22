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

// GetAllLinks returns all links.
func GetAllLinks(db *database.DB, page, limit, search, sortBy, sort string) (links []models.Link, total int64, err error) {
	// Total rows
	total_query := db.Model(&links)
	if search != "" {
		total_query.Where("url LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	total_query.Count(&total)

	var q = db.Scopes(database.Paginate(page, limit))
	if search != "" {
		q.Where("url LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if sort != "desc" {
		sort = "asc"
	}
	switch sortBy {
	case "url":
		q.Order("url " + sort)
	case "name":
		q.Order("name " + sort)
	case "expired_at":
		q.Order("expired_at " + sort)
	}

	if result := q.Find(&links); result.Error != nil {
		return links, total, result.Error
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
		Name:      link.Name,
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
func GetLinkFromID(db *database.DB, id string, expired bool) (link *models.Link, err error) {
	q := db
	if !expired {
		q.Where("expired_at >= ?", time.Now())
	}
	result := q.First(&link, "id = ?", id)
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
