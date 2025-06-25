package auth

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// Secret loaded from env to sign/verify tokens
var JWTSecret = []byte(os.Getenv("JWT_SECRET"))

// CustomClaims embeds the standard RegisteredClaims and adds own fields
type CustomClaims struct {
	UserID string `jsong:"sub"`
	jwt.RegisteredClaims
}