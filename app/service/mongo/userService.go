package mongo

import (
	// Beri alias 'mongoRepo' untuk repository
	mongoRepo "latihan2/app/repository/mongo"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// GetAllUsers adalah handler untuk (GET /users-m/mongo/)
func GetAllUsers(c *fiber.Ctx) error {
	// Parsing query params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "_id") // Default sort by _id
	order := strings.ToLower(c.Query("order", "asc"))
	search := c.Query("search", "")
	offset := (page - 1) * limit

	// Panggil Repository
	data, err := mongoRepo.GetUsersRepo(search, sortBy, order, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	total, err := mongoRepo.CountUsersRepo(search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Format response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
		"meta": fiber.Map{
			"page":   page,
			"limit":  limit,
			"total":  total,
			"pages":  (total + limit - 1) / limit,
			"sortBy": sortBy,
			"order":  order,
			"search": search,
		},
	})
}

// GetUsersByID adalah handler untuk (GET /users-m/mongo/:id/)
func GetUsersByID(c *fiber.Ctx) error {
	id := c.Params("id")

	// Panggil Repository
	data, err := mongoRepo.GetUserByID(id)
	if err != nil {
		if err == mongodriver.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}
