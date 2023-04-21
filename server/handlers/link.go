package handlers

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/fabienbellanger/go-url-shortener/server/db"
	models "github.com/fabienbellanger/go-url-shortener/server/models"
	"github.com/fabienbellanger/go-url-shortener/server/repositories"
	"github.com/fabienbellanger/go-url-shortener/server/utils"
	"github.com/fabienbellanger/goutils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type LinkImportError struct {
	Line int    `json:"line"`
	Err  string `json:"err"`
	Data string `json:"data"`
}

type LinksIdForm = []string

func newLinkImportError(index int, err string, line []string) LinkImportError {
	return LinkImportError{
		Line: index + 1,
		Err:  err,
		Data: fmt.Sprintf("\"%s\";\"%s\";\"%s\"", line[0], line[1], line[2]),
	}
}

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

// DeleteLinks delete selected links.
func DeleteLinks(db *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		links := new(LinksIdForm)
		if err := c.BodyParser(links); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
				Code:    fiber.StatusBadRequest,
				Message: "Bad Request",
			})
		}

		fmt.Printf("ids=%v\n", links)

		err := repositories.DeleteLinks(db, *links)
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
func UploadLink(db *db.DB, logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if form, err := c.MultipartForm(); err == nil {
			for _, files := range form.File {
				for _, file := range files {
					if file.Header["Content-Type"][0] == "text/csv" {
						fileName := fmt.Sprintf("./uploads/csv/%s.csv", uuid.New().String())

						// Save the files to disk
						if err := c.SaveFile(file, fileName); err != nil {
							return fiber.NewError(fiber.StatusInternalServerError, err.Error())
						}
						defer func() {
							// Remove file
							err = os.Remove(fileName)
							if err != nil {
								logger.Error(fmt.Sprintf("error when removing CSV upload file %v", fileName), zap.Error(err))
							}
						}()

						// Read file
						f, err := os.ReadFile(fileName)
						if err != nil {
							return fiber.NewError(fiber.StatusInternalServerError, err.Error())
						}

						// CSV parser
						reader := csv.NewReader(strings.NewReader(string(f)))
						reader.Comma = ';'
						records, err := reader.ReadAll()
						if err != nil {
							return fiber.NewError(fiber.StatusInternalServerError, err.Error())
						}

						// 100 lines max
						if len(records) > 100 {
							return c.Status(fiber.StatusBadRequest).JSON(utils.HTTPError{
								Code:    fiber.StatusBadRequest,
								Message: "Bad Request",
								Details: "Too many lines (max: 100)",
							})
						}

						// Save links in database
						var linesError []LinkImportError

						insertedLinks := 0
						for i, line := range records {
							// Skip CSV headers
							if i == 0 {
								continue
							}

							expiratedAt, err := goutils.SQLDatetimeToTime(line[2])
							if err != nil {
								linesError = append(linesError, newLinkImportError(i, "Invalid expiration date", line))
								continue
							}

							link := models.LinkForm{
								URL:       line[0],
								Name:      &line[1],
								ExpiredAt: expiratedAt,
							}

							_, err = repositories.CreateLink(db, &link)
							if err != nil {
								linesError = append(linesError, newLinkImportError(i, "Error during creation", line))
								continue
							}
							insertedLinks++
						}

						return c.Status(fiber.StatusOK).JSON(fiber.Map{
							"inserted_links": insertedLinks,
							"errors":         linesError,
						})
					}
				}
			}
		}

		return fiber.NewError(fiber.StatusInternalServerError, "invalid or empty CSV file")
	}
}
