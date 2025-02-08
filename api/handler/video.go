package handler

import (
	"math"
	"strconv"

	"github.com/TheAlpha16/typi/api/config"
	"github.com/TheAlpha16/typi/api/database"

	"github.com/gofiber/fiber/v2"
)

func GetVideos(c *fiber.Ctx) error {
	var err error

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		page = 1
	}

	perPage, err := strconv.Atoi(c.Query("per_page"))
	if err != nil {
		perPage = config.DEFAULT_PER_PAGE
	}

	if page < 1 {
		page = 1
	}

	if perPage < 1 || perPage > config.MAX_PER_PAGE {
		perPage = config.DEFAULT_PER_PAGE
	}

	offset := (page - 1) * perPage
	totalVideos, err := database.GetVideoCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failure",
			"message": "Error fetching videos",
		})
	}

	videos, err := database.GetVideos(offset, perPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "failure",
			"message": "Error fetching videos",
		})
	}
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
