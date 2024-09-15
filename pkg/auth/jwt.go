package auth

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type JWTClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(userID int64, secretKey string, expirationTime time.Duration) (string, error) {
	now := time.Now()
	expiresAt := now.Add(expirationTime)

	claims := &JWTClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  now.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secretKey))

	return tokenString, err
}

func ValidateToken(tokenString string, secretKey string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
