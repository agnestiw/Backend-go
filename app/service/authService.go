package service

import (
	"database/sql"
	"latihan2/app/model"
	"latihan2/app/repository"
	"latihan2/utils"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Request body tidak valid",
		})
	}
	if req.Username == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username dan password harus diisi",
		})
	}
	user, passwordHash, err := repository.GetUserByUsername(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Username atau password salah",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error database",
		})
	}
	if !utils.CheckPassword(req.Password, passwordHash) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Username atau password salah",
		})
	}
	token, err := utils.GenerateToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal generate token",
		})
	}
	response := model.LoginResponse{
		User:  *user,
		Token: token,
	}
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data":    response,
	})
}

func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	username := c.Locals("username").(string)
	role := c.Locals("role").(string)
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Profile berhasil diambil",
		"data": fiber.Map{
			"user_id":  userID,
			"username": username,
			"role":     role,
		},
	})
}