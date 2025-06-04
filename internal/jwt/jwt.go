package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrValidationJwt      = errors.New("unexpected signing method")
	ErrFailedToParseToken = errors.New("failed to parse token")
	ErrInvalidTokenClaims = errors.New("invalid token claims")
	ErrSubClaimNotFound   = errors.New("sub claim not found or not a string")
	ErrInvalidPassword    = errors.New("invalid password")
)

func CreateJwt(secret string, password string) (string, error) {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": string(passHash),
		"exp": time.Now().Add(time.Hour * 8).Unix(),
		"iat": time.Now().Unix(),
	})
	tokenString, err := claims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	fmt.Println(tokenString)
	return tokenString, nil
}

func ValidateJwt(tokenString, secret string, plainPassword string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrValidationJwt
		}
		return []byte(secret), nil
	})
	if err != nil {
		return false, ErrFailedToParseToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return false, ErrInvalidTokenClaims
	}

	hashedPassword, ok := claims["sub"].(string)
	if !ok {
		return false, ErrSubClaimNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		return false, ErrInvalidPassword
	}

	return true, nil
}
