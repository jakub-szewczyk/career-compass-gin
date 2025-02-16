package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Returns the health status of the service
//	@Tags			Health check
//	@Produce		json
//	@Success		200	{object}	models.HealthCheckResBody
//	@Router			/health-check [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}
