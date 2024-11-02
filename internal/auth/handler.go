package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/mwelwankuta/facebook-notes/pkg/config"
	"github.com/mwelwankuta/facebook-notes/pkg/models"
)

type AuthHandler struct {
	useCase           AuthUseCase
	OpenGraphClientID string
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(useCase AuthUseCase, clientId string) *AuthHandler {
	return &AuthHandler{
		useCase:           useCase,
		OpenGraphClientID: clientId,
	}
}

// AuthenticateUserHandler authenticates a user using the facebook code
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

// GetAllUsersHandler returns all users
func (a *AuthHandler) GetAllUsersHandler(c echo.Context) error {
	var dto = models.PaginateDto{
		Offset: "0",
		Limit:  "20",
	}

	if offset := c.QueryParam("page"); offset == "" {
		dto = models.PaginateDto{
			Offset: "0",
			Limit:  dto.Limit,
		}
	}
	if limit := c.QueryParam("limit"); limit == "" {
		dto = models.PaginateDto{
			Limit:  "20",
			Offset: dto.Offset,
		}
	}

	users, err := a.useCase.GetAllUsers(dto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return c.JSON(http.StatusOK, users)
}

// GetUserByIDHandler returns a user by ID
func (a *AuthHandler) GetUserByIDHandler(c echo.Context) error {
	userId := c.Param("id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "User ID not found"})
	}

	user, err := a.useCase.GetUserByID(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return c.JSON(http.StatusOK, user)
}

// LoginWithFacebook initiates the Facebook OAuth flow
func (a *AuthHandler) LoginWithFacebook(c echo.Context) error {
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=email", config.FbAuthURL, a.OpenGraphClientID, url.QueryEscape(config.RedirectURI))
	return c.String(http.StatusOK, authURL)
}
