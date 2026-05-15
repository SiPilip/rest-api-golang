package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func GenerateAccessToken(email string, userId int64, role string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatal("JWT_SECRET is not set in .env file")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"role": role,
		"exp": time.Now().Add(time.Second * 5).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func GenerateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func VerifyToken(tokenString string) (int64, string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return 0, "", errors.New("invalid token")
	}

	if !parsedToken.Valid {
		return 0, "", errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errors.New("invalid token")
	}

	// email := claims["email"].(string)
	userId := int64(claims["userId"].(float64))
	role := claims["role"].(string)
	return userId, role, nil
}