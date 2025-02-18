package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	common "github.com/jakub-szewczyk/career-compass-gin/utils"
)

// Profile godoc
//
//	@Summary		Get user profile
//	@Description	Retrieves and returns the profile information of the currently authenticated user
//
//	@Security		BearerAuth
//
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

	c.JSON(http.StatusOK, profileResBody)
}

// VeifyEmail godoc
//
//	@Summary		Verify user email
//	@Description	Confirms the email address of the currently authenticated user. This endpoint requires an email verification token sent to the user's registered email.
//
//	@Security		BearerAuth
//
//	@Tags			Profile
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.VerifyEmailReqBody	true	"Email verification data"
//	@Failure		400		{object}	models.Error
//	@Failure		500		{object}	models.Error
//	@Success		200		{object}	models.ProfileResBody
//	@Router			/profile/verify-email [patch]
func (h *Handler) VerifyEmail(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	uuid, err := common.ToUUID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var body models.VerifyEmailReqBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	t, err := h.queries.GetVerificationToken(h.ctx, uuid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if body.VerificationToken != t.Token {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid verification token",
		})
		return
	}

	if t.ExpiresAt.Time.Before(time.Now()) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "expired verification token",
		})
		return
	}

	user, err := h.queries.VerifyEmail(h.ctx, uuid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
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

	c.JSON(http.StatusOK, profileResBody)
}
