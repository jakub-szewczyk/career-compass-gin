package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	common "github.com/jakub-szewczyk/career-compass-gin/utils"
)

// Profile godoc
//
//	@Summary		Get user profile
//	@Description	Retrieves and returns the profile information of the currently authenticated user
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Success		200	{object}	models.ProfileResBody
//	@Router			/profile [get]
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

	profileResBody, err := models.NewProfileResBody(user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, *profileResBody)
}
