package jwt

import (
	"time"
	"uber-popug/pkg/types"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("supersecretkey")

func GenerateJWT(email, userID, role string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &types.JWTClaim{
		Email:    email,
		Username: userID,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}
