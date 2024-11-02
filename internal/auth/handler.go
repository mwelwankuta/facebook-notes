package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/mwelwankuta/facebook-notes/pkg/config"
)

type AuthHandler struct {
	useCase           AuthUseCase
	OpenGraphClientID string
}

func NewAuthHandler(useCase AuthUseCase, clientId string) *AuthHandler {
	return &AuthHandler{
		useCase:           useCase,
		OpenGraphClientID: clientId,
	}
}

func (a *AuthHandler) AuthenticateUserHandler(c echo.Context) error {
	code := c.QueryParam("code")
	if code == "" {
		return c.String(http.StatusBadRequest, "Code not found")
	}

	user, err := a.useCase.AuthenticateUser(code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return c.JSON(http.StatusOK, user)
}

// LoginWithFacebook initiates the Facebook OAuth flow
func (a *AuthHandler) LoginWithFacebook(c echo.Context) error {
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=email", config.FbAuthURL, a.OpenGraphClientID, url.QueryEscape(config.RedirectURI))
	return c.Redirect(http.StatusTemporaryRedirect, authURL)
}
