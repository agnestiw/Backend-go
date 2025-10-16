package utils

import (
	"fmt"
	"latihan2/app/model"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	return []byte(strings.TrimSpace(os.Getenv("JWT_SECRET_KEY")))
}

func GenerateToken(user model.User) (string, error) {
	secret := getJWTSecret()
	fmt.Printf("SECRET in GenerateToken: %q\n", secret)
	claims := model.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateToken(tokenString string) (*model.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*model.JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}
