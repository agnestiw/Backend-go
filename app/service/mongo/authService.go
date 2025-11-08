package mongo

import (
	"errors"
	"fmt"
	"latihan2/app/model" 
	mongoModel "latihan2/app/model/mongo"
	mongoRepo "latihan2/app/repository/mongo"
	"latihan2/utils"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// LoginMongo godoc
// @Summary Login user
// @Description Login untuk mendapatkan JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Data login"
// @Success 200 {object} mongoModel.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/mg/login [post]
func LoginMongo(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Request body tidak valid",
		})
	}
	user, err := mongoRepo.GetUserByUsername(req.Username)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Username atau password salah",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Username atau password salah",
		})
	}

	fmt.Printf(">>> [DEBUG] Memanggil GenerateMongoToken dengan User: %+v", user)
	token, err := utils.GenerateMongoToken(user.ID, user.Username, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuat token",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"data": mongoModel.LoginResponse{
			User: mongoModel.UserMongo{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				Role:      user.Role,
				CreatedAt: user.CreatedAt,
			},
			Token: token,
		},
	})
}

func GetProfileMongo(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "UserID tidak ditemukan dari token",
		})
	}

	user, err := mongoRepo.GetUserByID(userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "User profile tidak ditemukan",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    user,
	})
}

func GenerateMongoToken(user *mongoModel.User) (string, error) {
	claims := jwt.MapClaims{
		"userID":   user.ID.Hex(),
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		secretKey = "your-default-secret-key"
	}

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", errors.New("failed to create token")
	}

	return tokenString, nil
}
