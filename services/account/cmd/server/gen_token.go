package server

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateToken() {
	secret := []byte("your-secret-key")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "admin",
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		panic(err)
	}

	fmt.Println("JWT:", tokenString)
}
