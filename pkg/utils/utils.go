package utils

import (
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/mwelwankuta/facebook-notes/pkg/models"
)

var validate = validator.New()

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(message string) ErrorResponse {
	return ErrorResponse{Error: message}
}

// GenerateJwtToken generates a jwt token
// Middleware exists to automatically read the token from the request and verify it
func GenerateJwtToken(secret string, user models.User, facebookToken string) (string, error) {
	claims := &JwtCustomClaims{
		user.ID,
		user.Role,
		facebookToken,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func EndpointNotImplemented(c echo.Context) error {
	return c.String(404, "Endpoint not implemented")
}

// GetPaginationFromQuery extracts pagination parameters from the request query
func GetPaginationFromQuery(c echo.Context) models.PaginateDto {
	// Default values
	defaultPage := 0
	defaultLimit := 20

	queryDto := models.GetPaginationFromQueryDto{
		Offset: c.QueryParam("page"),
		Limit:  c.QueryParam("limit"),
	}

	// Convert offset
	pageNum, err := strconv.Atoi(queryDto.Offset)
	if err != nil || queryDto.Offset == "" {
		pageNum = defaultPage
	}

	// Convert limit
	limitNum, err := strconv.Atoi(queryDto.Limit)
	if err != nil || queryDto.Limit == "" {
		limitNum = defaultLimit
	}

	return models.PaginateDto{
		Offset: pageNum,
		Limit:  limitNum,
	}
}

// Validate validates a struct using validator.v10 package
func Validate(s interface{}) error {
	return validate.Struct(s)
}
