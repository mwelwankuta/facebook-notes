package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mwelwankuta/facebook-notes/pkg/models"
)

// GenerateJwtToken generates a jwt token
// Middleware exists to automatically read the token from the request and verify it
func GenerateJwtToken(secret string, user models.User, facebookToken string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":            time.Now().Add(time.Hour * 72).Unix(), // 3 days
		"iat":            time.Now().Unix(),
		"user":           user,
		"facebook_token": facebookToken,
	})

	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func EndpointNotImplemented(c echo.Context) error {
	return c.String(404, "Endpoint not implemented")
}
