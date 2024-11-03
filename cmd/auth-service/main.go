package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

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

	// jwtMiddlware := echojwt.WithConfig(echojwt.Config{
	// 	SigningKey:  []byte(cfg.JwtSecret),
	// 	TokenLookup: "header:Authorization",
	// })

	// authentication
	e.POST("/api/auth/login/callback", authHandler.AuthenticateUserHandler)
	// endpoint called from client
	e.GET("/api/auth/login", authHandler.LoginWithFacebook)

	// users
	e.GET("/api/auth/users", authHandler.GetAllUsersHandler)
	e.GET("/api/auth/users/:id", authHandler.GetUserByIDHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}
