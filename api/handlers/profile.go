package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/jakub-szewczyk/career-compass-gin/utils"
)

func (h *Handler) Profile(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	uuid, err := common.ToUUID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := h.queries.GetUserById(h.ctx, uuid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}
