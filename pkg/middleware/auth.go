package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/mwelwankuta/facebook-notes/pkg/utils"
)

func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := utils.GetUserFromContext(c)
			if err != nil {
				return echo.NewHTTPError(401, "unauthorized")
			}

			for _, role := range roles {
				if user.Role == role {
					return next(c)
				}
			}

			return echo.NewHTTPError(403, "forbidden")
		}
	}
}
