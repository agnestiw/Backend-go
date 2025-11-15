package middleware

import (
	"fmt"
	"latihan2/app/repository/mongo"
	"latihan2/utils"
	"log"

	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	colorBlue  = "\033[34m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

// untuk PostgreSQL
func AuthRequired() fiber.Handler {
    // fmt.Println(">>> [AuthRequired] Middleware dijalankan")

	return func(c *fiber.Ctx) error {
		path := c.Path()
		log.Printf("%s[AuthRequired] Middleware aktif untuk path: %s%s", colorBlue, path, colorReset)

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			log.Printf("%s[AuthRequired] ❌ Tidak ada Authorization header di path: %s%s", colorRed, path, colorReset)
			return c.Status(401).JSON(fiber.Map{
				"error": "Token akses diperlukan",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Printf("%s[AuthRequired] ⚠️ Format token salah di path: %s | Header: %s%s", colorRed, path, authHeader, colorReset)
			return c.Status(401).JSON(fiber.Map{
				"error": "Format token tidak valid",
			})
		}

		claims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			log.Printf("%s[AuthRequired] ❌ Token tidak valid atau expired di path: %s | Error: %v%s", colorRed, path, err, colorReset)
			return c.Status(401).JSON(fiber.Map{
				"error": "Token tidak valid atau expired",
			})
		}

		idStr := fmt.Sprintf("%v", claims.UserID)
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("%s[AuthRequired] ⚠️ User ID di token tidak valid di path: %s%s", colorRed, path, colorReset)
			return c.Status(401).JSON(fiber.Map{
				"error": "User ID di token tidak valid",
			})
		}

		log.Printf("%s[AuthRequired] ✅ Token valid. UserID: %d | Username: %v | Role: %v | Path: %s%s",
			colorGreen, id, claims.Username, claims.Role, path, colorReset)

		// simpan data user ke context
		c.Locals("user_id", id)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func AdminOnly() fiber.Handler {
    return func(c *fiber.Ctx) error {
        role, ok := c.Locals("role").(string)
        if !ok || role != "admin" {
            return c.Status(403).JSON(fiber.Map{
                "error": "Akses ditolak. Hanya admin yang diizinkan",
            })
        }
        return c.Next()
    }
}


// untuk MongoDB
func AuthRequiredMongo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token akses diperlukan",
			})
		}

		// Pisahkan token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Format token tidak valid, gunakan 'Bearer <token>'",
			})
		}

		// ✅ Panggil fungsi dari jwtmongo.go
		claims, err := utils.ValidateMongoToken(tokenParts[1])
		if err != nil || claims == nil {
			log.Printf("[DEBUG-AuthMongo] Error Validasi Token: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token tidak valid atau expired",
			})
		}

		// Debug info
		// log.Printf("[DEBUG-AuthMongo] Claims Diterima: UserID='%s', Username='%s', Role='%s'",
		// 	claims.UserID, claims.Username, claims.Role)

		if claims.UserID == "" || claims.Role == "" {
			log.Printf("[DEBUG-AuthMongo] Validasi Gagal: UserID atau Role kosong!")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Informasi user di token tidak lengkap",
			})
		}

		// Simpan ke context agar handler bisa akses
		c.Locals("userID", claims.UserID)
		c.Locals("role", claims.Role)
		c.Locals("username", claims.Username)

		return c.Next()
	}
}


func FileOwnerOrAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fileID := c.Params("id")
		role := c.Locals("role").(string)
		loggedInUserID := c.Locals("userID").(string)

		file, err := mongo.FindFileByID(fileID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "File tidak ditemukan"})
		}

		if role != "admin" && file.OwnerID.Hex() != loggedInUserID {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak"})
		}

		return c.Next()
	}
}
