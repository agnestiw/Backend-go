package middleware

import (
	"fmt"
	"latihan2/utils"
	
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	
)

func AuthRequired() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Token akses diperlukan",
            })
        }
        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            return c.Status(401).JSON(fiber.Map{
                "error": "Format token tidak valid",
            })
        }
        claims, err := utils.ValidateToken(tokenParts[1])
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "Token tidak valid atau expired",
            })
        }

        idStr := fmt.Sprintf("%v", claims.UserID)
        id, err := strconv.Atoi(idStr)
        if err != nil {
            return c.Status(401).JSON(fiber.Map{
                "error": "User ID di token tidak valid",
            })
        }

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


// func AuthMiddlewareAdmin() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		authHeader := c.Get("Authorization")
// 		if authHeader == "" {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Header otorisasi tidak ditemukan"})
// 		}

// 		parts := strings.Split(authHeader, " ")
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Format token tidak valid"})
// 		}
// 		tokenString := parts[1]
        
//         log.Println("--- DEBUGGING AUTH MIDDLEWARE ---")
//         log.Println("Authorization Header:", c.Get("Authorization"))

// 		claims := jwt.MapClaims{}
// 		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fiber.NewError(fiber.StatusUnauthorized, "Metode signing tidak terduga")
// 			}
// 			// DIUBAH: Menggunakan "JWT_SECRET_KEY" agar konsisten
// 			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
// 		})

// 		if err != nil || !token.Valid {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid atau kedaluwarsa"})
// 		}

// 		role, ok := claims["role"].(string)
// 		if !ok || role != "admin" {
// 			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak. Hanya admin yang diizinkan"})
// 		}
        
// 		// DIUBAH: Menggunakan "user_id" agar konsisten
// 		userID, ok := claims["user_id"].(float64)
// 		if !ok {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "UserID tidak ditemukan di dalam token"})
// 		}

// 		// DIUBAH: Menggunakan "user_id" agar konsisten
// 		c.Locals("user_id", int(userID))
// 		c.Locals("role", role)

// 		return c.Next()
// 	}
// }