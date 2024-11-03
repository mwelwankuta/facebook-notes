package main

import (
	"fmt"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/mwelwankuta/facebook-notes/internal/summaries"
	"github.com/mwelwankuta/facebook-notes/pkg/config"
	"github.com/mwelwankuta/facebook-notes/pkg/db"
	"github.com/mwelwankuta/facebook-notes/pkg/utils"
)

func main() {
	cfg, err := config.LoadConfig("config/auth-config.yaml")
	if err != nil {
		panic("Could not load config file")
	}

	database := db.InitializeDatabase(cfg.Database)

	summariesRepository := summaries.NewSummariesRepository(database)
	summariesUseCase := summaries.NewSummariesUseCase(*summariesRepository, *cfg)
	summaries.NewSummariesHandler(*summariesUseCase, cfg.OpenGraphClientID)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(cfg.JwtSecret),
		TokenLookup: "header:Authorization",
	})

	e.GET("/api/summaries", utils.EndpointNotImplemented)
	e.GET("/api/summaries/requests", utils.EndpointNotImplemented)
	e.POST("/api/summaries/requests", utils.EndpointNotImplemented, jwtMiddleware)

	e.GET("/api/summaries/:id", utils.EndpointNotImplemented)
	e.POST("/api/summaries/:id/rate", utils.EndpointNotImplemented, jwtMiddleware)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}
