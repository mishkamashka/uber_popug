package types

import "github.com/dgrijalva/jwt-go"

type JWTClaim struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.StandardClaims
}
