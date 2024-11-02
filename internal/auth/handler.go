package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	useCase AuthUseCase
}

func NewAuthHandler(useCase AuthUseCase) *AuthHandler {
	return &AuthHandler{useCase: useCase}
}

func (a *AuthHandler) AuthenticateUserHandler(c echo.Context) error {
	userDto := FacebookUser{}
	if err := c.Bind(userDto); err != nil {
		return c.JSON(http.StatusOK, "Invalid request payload")
	}

	user, err := a.useCase.AuthenticateUser(userDto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal server error"})
	}

	return c.JSON(http.StatusOK, user)
}
