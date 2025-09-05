package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	"github.com/jakub-szewczyk/career-compass-gin/utils"
)

func (h *Handler) Resumes(c *gin.Context) {}

func (h *Handler) Resume(c *gin.Context) {}

// CreateResume godoc
//
//	@Summary		Generate a new resume
//	@Description	Generates a new resume. If no title is provided, a default one will be used.
//
//	@Security		BearerAuth
//	@Tags			Resume
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.CreateResumeReqBody	false	"Resume details"
//	@Failure		400		{object}	models.Error
//	@Failure		500		{object}	models.Error
//	@Success		201		{object}	models.CreateResumeResBody
//	@Router			/resumes [post]
func (h *Handler) CreateResume(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	uuid, err := utils.ToUUID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var body models.CreateResumeReqBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	resume, err := h.queries.CreateResume(h.ctx, db.CreateResumeParams{
		UserID: uuid,
		Title:  body.Title,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resBody := models.NewCreateResumeResBody(resume)

	c.JSON(http.StatusCreated, resBody)
}

func (h *Handler) UpdateResume(c *gin.Context) {}

// DeleteResume godoc
//
//	@Summary		Delete a resume
//	@Description	Deletes an existing resume
//
//	@Security		BearerAuth
//
//	@Tags			Resume
//	@Accept			json
//	@Produce		json
//	@Param			resumeId	path		string	true	"Resume uuid"
//	@Success		200			{object}	models.DeleteResumeResBody
//	@Failure		500			{object}	models.Error
//	@Router			/resumes/{resumeId} [delete]
func (h *Handler) DeleteResume(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	uuid, err := utils.ToUUID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resumeId, err := utils.ToUUID(c.Param("resumeId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resume, err := h.queries.DeleteResume(h.ctx, db.DeleteResumeParams{
		ID:     resumeId,
		UserID: uuid,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resBody := models.NewDeleteResumeResBody(resume)

	c.JSON(http.StatusOK, resBody)
}
