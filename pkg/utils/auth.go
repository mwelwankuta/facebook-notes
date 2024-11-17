package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mwelwankuta/facebook-notes/pkg/models"
)

type JwtCustomClaims struct {
	ID          string `json:"id"`
	Role        string `json:"role"`
	AccessToken string `json:"access_token"`
	jwt.RegisteredClaims
}

func GetUserFromContext(c echo.Context) (models.User, error) {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*JwtCustomClaims)

	if claims == nil {
		return models.User{}, errors.New("invalid token claims")
	}

	return models.User{
		ID:   claims.ID,
		Role: claims.Role,
	}, nil
}
