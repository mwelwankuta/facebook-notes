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
	db.InitializeDatabase(cfg.Database)

	authRepository := auth.NewAuthRepository()
	authUseCase := auth.NewAuthUseCase(*authRepository)
	authHandler := auth.NewAuthHandler(*authUseCase)

	e := echo.New()
	e.POST("/api/authenticate", authHandler.AuthenticateUserHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}
