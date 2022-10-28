package handlers

import (
	"github.com/fabienbellanger/go-url-shortener/server/db"
	models "github.com/fabienbellanger/go-url-shortener/server/models"
	"github.com/fabienbellanger/go-url-shortener/server/repositories"
	"github.com/fabienbellanger/go-url-shortener/server/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// LinksList returns the list of all links
func LinksList(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page := c.Query("page")
		limit := c.Query("limit")
		search := c.Query("s")
		sortBy := c.Query("sort-by")
		sort := c.Query("sort") // asc or desc

		links, total, err := repositories.GetAllLinks(db, page, limit, search, sortBy, sort)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Error when getting all links")
		}
		return c.JSON(fiber.Map{
			"total": total,
			"links": links,
		})
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
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(newLink)
	}
}

func UpdateLink(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		if len(id) != 8 {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		link, err := repositories.GetLinkFromID(db, id, true)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "No link found")
		}

		linkForm := new(models.LinkForm)
		if err := c.BodyParser(linkForm); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		// Data validation
		// ---------------
		createErrors := utils.ValidateStruct(*linkForm)
		if createErrors != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
				Details: createErrors,
			})
		}

		link.URL = linkForm.URL
		link.Name = linkForm.Name
		link.ExpiredAt = linkForm.ExpiredAt

		err = repositories.UpdateLink(db, link)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// DeleteLink delete the specified link.
func DeleteLink(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		if len(id) != 8 {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		err := repositories.DeleteLink(db, id)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// RedirectURL redirects to original URL if URL exists, else return 404 HTTP code.
func RedirectURL(db *db.DB, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if id == "" {
			c.SendStatus(fiber.StatusNotFound)
		}

		link, err := repositories.GetLinkFromID(db, id, false)
		if err != nil {
			return c.SendStatus(fiber.StatusNotFound)
		}
		return c.Redirect(link.URL)
	}
}

// UploadLink creates links from CSV file.
func UploadLink(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("OK")
	}
}
