package handlers

import (
	"github.com/fabienbellanger/go-url-shortener/db"
	models "github.com/fabienbellanger/go-url-shortener/models"
	"github.com/fabienbellanger/go-url-shortener/repositories"
	"github.com/fabienbellanger/go-url-shortener/utils"
	"github.com/gofiber/fiber/v2"
)

// LinksList returns the list of all links
// TODO: Paginate!!
func LinksList(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.Query("page")
		limit := c.Query("limit")

		links, err := repositories.GetAllLinks(db, page, limit)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when getting all links")
		}
		return c.JSON(links)
	}
}

// CreateLink adds a new link.
func CreateLink(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		link := new(models.LinkForm)
		if err := c.BodyParser(link); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		// Data validation
		// ---------------
		createErrors := utils.ValidateStruct(*link)
		if createErrors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
				Details: createErrors,
			})
		}

		newLink, err := repositories.CreateLink(db, link)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when creating link")
		}

		return c.JSON(newLink)
	}
}

// RedirectURL redirects to original URL if URL exists, else return 404 HTTP code.
func RedirectURL(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			c.SendStatus(fiber.StatusNotFound)
		}

		link, err := repositories.GetLinkFromID(db, id)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Redirect(link.URL)
	}
}
