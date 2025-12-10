package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jonesrussell/gosources/internal/logger"
	"github.com/jonesrussell/gosources/internal/models"
	"github.com/jonesrussell/gosources/internal/repository"
)

// GlobalSelectorsHandler handles HTTP requests for global selectors
type GlobalSelectorsHandler struct {
	repo   *repository.GlobalSelectorsRepository
	logger logger.Logger
}

// NewGlobalSelectorsHandler creates a new GlobalSelectorsHandler
func NewGlobalSelectorsHandler(repo *repository.GlobalSelectorsRepository, log logger.Logger) *GlobalSelectorsHandler {
	return &GlobalSelectorsHandler{
		repo:   repo,
		logger: log,
	}
}

// Get retrieves the global selectors configuration
func (h *GlobalSelectorsHandler) Get(c *gin.Context) {
	gs, err := h.repo.Get(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to get global selectors",
			logger.Error(err),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": "Global selectors not found"})
		return
	}

	c.JSON(http.StatusOK, gs)
}

// Update updates the global selectors configuration
func (h *GlobalSelectorsHandler) Update(c *gin.Context) {
	var selectors models.SelectorConfig
	if err := c.ShouldBindJSON(&selectors); err != nil {
		h.logger.Debug("Invalid request body",
			logger.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), &selectors); err != nil {
		h.logger.Error("Failed to update global selectors",
			logger.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update global selectors"})
		return
	}

	h.logger.Info("Global selectors updated")

	// Return the updated configuration
	gs, err := h.repo.Get(c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to retrieve updated global selectors",
			logger.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated global selectors"})
		return
	}

	c.JSON(http.StatusOK, gs)
}
