package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jonesrussell/gosources/internal/handlers"
	"github.com/jonesrussell/gosources/internal/logger"
	"github.com/jonesrussell/gosources/internal/repository"
)

func NewRouter(db *repository.SourceRepository, log logger.Logger) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(ginLogger(log))
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1
	v1 := router.Group("/api/v1")
	{
		sourceHandler := handlers.NewSourceHandler(db, log)

		// Sources endpoints
		sources := v1.Group("/sources")
		{
			sources.POST("", sourceHandler.Create)
			sources.GET("", sourceHandler.List)
			sources.GET("/:id", sourceHandler.GetByID)
			sources.PUT("/:id", sourceHandler.Update)
			sources.DELETE("/:id", sourceHandler.Delete)
		}

		// Cities endpoint for gopost integration
		v1.GET("/cities", sourceHandler.GetCities)
	}

	return router
}

func ginLogger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		log.Info("HTTP request",
			logger.String("method", method),
			logger.String("path", path),
			logger.Int("status_code", statusCode),
			logger.String("client_ip", c.ClientIP()),
			logger.Duration("duration", duration),
		)
	}
}
