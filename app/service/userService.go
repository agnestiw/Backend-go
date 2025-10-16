package service

import (
	"database/sql"
	"latihan2/app/model"
	"latihan2/app/repository"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetUsersService(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	sortBy := c.Query("sortBy", "id")
	order := c.Query("order", "asc")
	search := c.Query("search", "")
	offset := (page - 1) * limit
	
	sortByWhitelist := map[string]bool{
		"id": true, 
		"name": true,
		"email": true, 
		"created_at": true}
	if !sortByWhitelist[sortBy] {
		sortBy = "id"
	}
	if strings.ToLower(order) != "desc" {
		order = "asc"
	}
	users, err := repository.GetUsersRepo(search, sortBy, order, c.Locals("role").(string), limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	total, err := repository.CountUsersRepo(search)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count users"})
	}
	response := model.UserResponse{
		Data: users,
		Meta: model.MetaInfo{
			Page:   page,
			Limit:  limit,
			Total:  total,
			Pages:  (total + limit - 1) / limit,
			SortBy: sortBy,
			Order:  order,
			Search: search,
		},
	}
	return c.JSON(response)
}

func SoftDeleteUserService(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
	return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	role := c.Locals("role")
	if role != "admin" {
		return c.Status(403).JSON(fiber.Map{"error": "Forbidden: only admin can delete"})
	}

	err = repository.SoftDeleteUserRepo(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.JSON(fiber.Map{"message": "User soft deleted successfully"})
}

func GetUserByIDService(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
    }

    role := c.Locals("role").(string)

    user, err := repository.GetUserByID(id, role)
    if err != nil {
        if err == sql.ErrNoRows {
            return c.Status(404).JSON(fiber.Map{"error": "User not found"})
        }
        return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user"})
    }

    return c.JSON(fiber.Map{
        "success": true,
        "data":    user,
        "message": "User fetched successfully",
    })
}
