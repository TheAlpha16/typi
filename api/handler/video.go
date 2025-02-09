package handler

import (
	"math"
	"strconv"

	"github.com/TheAlpha16/typi/api/config"
	"github.com/TheAlpha16/typi/api/database"

	"github.com/gofiber/fiber/v2"
)

// GetVideos retrieves paginated video data from the database
func GetVideos(c *fiber.Ctx) error {
	var err error

	// page is the current page number (defaults to 1 if invalid)
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		page = 1
	}

	// perPage is the maximum items in one result page
	perPage, err := strconv.Atoi(c.Query("per_page"))
	if err != nil {
		perPage = config.DEFAULT_PER_PAGE
	}

	// ensure page and perPage are within allowed limits
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > config.MAX_PER_PAGE {
		perPage = config.DEFAULT_PER_PAGE
	}

	// offset is the starting index for database retrieval
	offset := (page - 1) * perPage

	// totalVideos holds the total number of videos in the database
	totalVideos, err := database.GetVideoCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failure",
			"message": "Error fetching videos",
		})
	}

	// videos slice contains the fetched video entries
	videos, err := database.GetVideos(offset, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failure",
			"message": "Error fetching videos",
		})
	}

	// totalPages is the total number of pages based on the perPage limit
	totalPages := int(math.Ceil(float64(totalVideos) / float64(perPage)))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":       "success",
		"total_videos": totalVideos,
		"page":         page,
		"per_page":     perPage,
		"total_pages":  totalPages,
		"videos":       videos,
	})
}
