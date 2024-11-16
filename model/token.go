package model

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	AcceessToken           string
	RefreshToken           string
	RefreshTokenExpiration time.Time
}

type CustomClaim struct {
	jwt.RegisteredClaims
	Name    string
	Email   string
	IsAdmin bool
	ID      uint
}
