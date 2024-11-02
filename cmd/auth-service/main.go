package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/mwelwankuta/facebook-notes/internal/auth"
	"github.com/mwelwankuta/facebook-notes/pkg/config"
	"github.com/mwelwankuta/facebook-notes/pkg/db"
)

func main() {
	cfg := config.LoadConfig("config/auth-config.yaml")
	db := db.InitializeDatabase(cfg.Database)

	authRepository := auth.NewAuthRepository(db)
	authUseCase := auth.NewAuthUseCase(*authRepository, cfg.OpenGraphClientID, cfg.OpenGraphClientSecret)
	authHandler := auth.NewAuthHandler(*authUseCase, cfg.OpenGraphClientID)

	e := echo.New()
	e.POST("/api/facebook/login/callback", authHandler.AuthenticateUserHandler)
	e.GET("/api/facebook/login", authHandler.LoginWithFacebook)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}
