package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/mwelwankuta/facebook-notes/pkg/config"
	"github.com/mwelwankuta/facebook-notes/pkg/utils"
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
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Code not found"))
	}

	user, err := a.useCase.AuthenticateUser(code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal server error"))
	}

	return c.JSON(http.StatusOK, user)
}

// GetAllUsersHandler returns all users
func (a *AuthHandler) GetAllUsersHandler(c echo.Context) error {
	dto := utils.GetPaginationFromQuery(c)
	users, err := a.useCase.GetAllUsers(dto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal server error"))
	}

	return c.JSON(http.StatusOK, users)
}

// GetUserByIDHandler returns a user by ID
func (a *AuthHandler) GetUserByIDHandler(c echo.Context) error {
	userId := c.Param("id")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("User ID not found"))
	}

	user, err := a.useCase.GetUserByID(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal server error"))
	}

	return c.JSON(http.StatusOK, user)
}

// LoginWithFacebook initiates the Facebook OAuth flow
func (a *AuthHandler) LoginWithFacebook(c echo.Context) error {
	authURL := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&scope=email", config.FbAuthURL, a.OpenGraphClientID, url.QueryEscape(config.RedirectURI))
	return c.String(http.StatusOK, authURL)
}

func (a *AuthHandler) GetCurrentUser(c echo.Context) error {
	user, err := utils.GetUserFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Unauthorized"))
	}
	return c.JSON(http.StatusOK, user)
}

func (a *AuthHandler) UpdateUserRole(c echo.Context) error {
	userId := c.Param("id")
	var req UpdateRoleRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request"))
	}

	if err := utils.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(err.Error()))
	}

	user, err := a.useCase.UpdateUserRole(userId, req.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, user)
}

func (a *AuthHandler) UpdateUserStatus(c echo.Context) error {
	userId := c.Param("id")
	var req UpdateStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse("Invalid request"))
	}

	if err := utils.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewErrorResponse(err.Error()))
	}

	user, err := a.useCase.UpdateUserStatus(userId, req.IsActive)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.NewErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, user)
}
