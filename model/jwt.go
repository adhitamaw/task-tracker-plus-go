package model

import "github.com/golang-jwt/jwt"

var JwtKey = []byte("secret-key")

type Claims struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}
