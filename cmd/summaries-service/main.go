package main

import (
	"fmt"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/mwelwankuta/facebook-notes/internal/summaries"
	"github.com/mwelwankuta/facebook-notes/pkg/config"
	"github.com/mwelwankuta/facebook-notes/pkg/db"
)

func main() {
	cfg, err := config.LoadConfig("config/summaries-config.yaml")
	if err != nil {
		panic("Could not load config file")
	}

	database := db.InitializeDatabase(cfg.Database)

	summariesRepository := summaries.NewSummariesRepository(database)
	summariesUseCase := summaries.NewSummariesUseCase(*summariesRepository, *cfg)
	summariesHandler := summaries.NewSummariesHandler(*summariesUseCase, cfg.OpenGraphClientID)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(cfg.JwtSecret),
		TokenLookup: "header:Authorization",
	})

	// Protected routes requiring authentication
	protected := e.Group("")
	protected.Use(jwtMiddleware)

	// User routes
	protected.POST("/api/summaries/requests", summariesHandler.CreateSummaryRequestHandler)
	protected.POST("/api/summaries/:id/rate", summariesHandler.RateSummaryHandler)

	// Moderator routes
	protected.POST("/api/summaries/:id/moderate", summariesHandler.ModerateSummaryHandler)
	protected.PUT("/api/summaries/:id/edit", summariesHandler.EditSummaryHandler)
	protected.POST("/api/summaries/:id/resources", summariesHandler.AddResourceLinkHandler)
	protected.DELETE("/api/summaries/:id/resources/:linkId", summariesHandler.RemoveResourceLinkHandler)

	// Public routes
	e.GET("/api/summaries", summariesHandler.GetAllSummariesHandler)
	e.GET("/api/summaries/requests", summariesHandler.GetAllRequestsHandler)
	e.GET("/api/summaries/:id", summariesHandler.GetSummaryByIDHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}
