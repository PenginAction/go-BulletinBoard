package utils

import (
	"log"
	"time"

	"github.com/PenginAction/go-BulletinBoard/config"
	"github.com/PenginAction/go-BulletinBoard/dto"
	"github.com/golang-jwt/jwt/v5"
)

func CreateValidToken(userID uint) (string, error) {
	cfg, err := config.LoadConfig("../")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	claims := &dto.JwtCustomClaims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.SECRET))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
