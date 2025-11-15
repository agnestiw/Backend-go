package utils

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtSecretKeyMongo = []byte("16824af3-6b8e-4c3d-9f1e-2c4b5e6f7g8h")

type MongoClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateMongoToken(userID primitive.ObjectID, username, role string) (string, error) {
	claims := MongoClaims{
		UserID:   userID.Hex(),
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKeyMongo)
	if err != nil {
		log.Printf("[ERROR-GenerateMongoToken] Gagal membuat token: %v", err)
		return "", err
	}
	log.Printf("[DEBUG-GenerateMongoToken] Token dibuat untuk UserID=%s, Username=%s, Role=%s",
		claims.UserID, claims.Username, claims.Role)

	return tokenString, nil
}

func ValidateMongoToken(tokenString string) (*MongoClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MongoClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKeyMongo, nil
	})
	if err != nil {
		log.Printf("[DEBUG-JWT-ValidateMongo] Token gagal di-parse: %v", err)
		return nil, err
	}
	claims, ok := token.Claims.(*MongoClaims)
	if !ok || !token.Valid {
		log.Printf("[DEBUG-JWT-ValidateMongo] Token tidak valid / gagal casting claims")
		return nil, errors.New("token tidak valid")
	}
	// log.Printf("[DEBUG-JWT-ValidateMongo] Parsed: UserID='%s', Username='%s', Role='%s'",
	// 	claims.UserID, claims.Username, claims.Role)
	if claims.UserID == "" || claims.Role == "" {
		// log.Printf("[DEBUG-JWT-ValidateMongo] Klaim kosong: UserID='%s', Role='%s'", claims.UserID, claims.Role)
		return nil, errors.New("klaim token tidak lengkap")
	}

	return claims, nil
}
