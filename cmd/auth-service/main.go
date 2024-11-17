package main

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	customMiddleware "github.com/mwelwankuta/facebook-notes/pkg/middleware"
	"github.com/mwelwankuta/facebook-notes/pkg/models"
	"github.com/mwelwankuta/facebook-notes/pkg/utils"

	"github.com/mwelwankuta/facebook-notes/internal/auth"
	"github.com/mwelwankuta/facebook-notes/pkg/config"
	"github.com/mwelwankuta/facebook-notes/pkg/db"
)

func main() {
	cfg, err := config.LoadConfig("config/auth-config.yaml")
	if err != nil {
		panic("Could not load config file")
	}

	database := db.InitializeDatabase(cfg.Database)

	authRepository := auth.NewAuthRepository(database)
	authUseCase := auth.NewAuthUseCase(*authRepository, *cfg)
	authHandler := auth.NewAuthHandler(*authUseCase, cfg.OpenGraphClientID)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(utils.JwtCustomClaims)
		},
		SigningKey: []byte(cfg.JwtSecret),
	}

	// Public routes
	e.POST("/api/auth/login/callback", authHandler.AuthenticateUserHandler)
	e.GET("/api/auth/login", authHandler.LoginWithFacebook)

	// Protected routes
	api := e.Group("/api")
	api.Use(echojwt.WithConfig(config))

	// User routes
	api.GET("/auth/users/me", authHandler.GetCurrentUser)
	api.GET("/auth/users", authHandler.GetAllUsersHandler)
	api.GET("/auth/users/:id", authHandler.GetUserByIDHandler)

	// Moderator routes
	moderator := api.Group("/admin")
	moderator.Use(customMiddleware.RequireRole(models.RoleModerator, models.RoleAdmin))
	moderator.PUT("/users/:id/role", authHandler.UpdateUserRole)
	moderator.PUT("/users/:id/status", authHandler.UpdateUserStatus)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}
