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
// @Summary		Generate a new resume
// @Description	Generates a new resume. If no title is provided, a default one will be used.
// @Tags			Resume
// @Accept			json
// @Produce		json
// @Param			body	body		models.CreateResumeReqBody	false	"Resume data"
// @Success		201		{object}	models.CreateResumeResBody
// @Failure		400		{object}	models.Error
// @Failure		500		{object}	models.Error
// @Router			/resumes [post]
// @Security		BearerAuth
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

func (h *Handler) DeleteResume(c *gin.Context) {}
