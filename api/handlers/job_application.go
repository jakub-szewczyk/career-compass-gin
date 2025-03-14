package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	common "github.com/jakub-szewczyk/career-compass-gin/utils"
)

// TODO
func (h *Handler) JobApplications(c *gin.Context) {}

// JobApplication godoc
//
//	@Summary		Retrieve job application details
//	@Description	Fetches the details of a specific job application by its id
//	@Tags			Job application
//	@Accept			json
//	@Produce		json
//	@Param			jobApplicationId	path		string	true	"Job application uuid"
//	@Failure		400					{object}	models.Error
//	@Failure		404					{object}	models.Error
//	@Failure		500					{object}	models.Error
//	@Success		200					{object}	models.JobApplicationResBody
//	@Router			/job-applications/{jobApplicationId} [get]
func (h *Handler) JobApplication(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	uuid, err := common.ToUUID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	jobApplicationId, err := common.ToUUID(c.Param("jobApplicationId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	jobApplication, err := h.queries.GetJobApplication(h.ctx, db.GetJobApplicationParams{
		ID:     jobApplicationId,
		UserID: uuid,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	jobApplicationResBody := models.NewJobApplicationResBody(jobApplication)

	c.JSON(http.StatusOK, jobApplicationResBody)
}

// CreateJobApplication godoc
//
//	@Summary		Submit a new job application
//	@Description	Processes and creates a new job application with the provided data
//	@Tags			Job application
//	@Accept			json
//	@Produce		json
//	@Param			body	body		models.CreateJobApplicationReqBody	true	"Job application details"
//	@Failure		400		{object}	models.Error
//	@Failure		500		{object}	models.Error
//	@Success		201		{object}	models.CreateJobApplicationResBody
//	@Router			/job-applications [post]
func (h *Handler) CreateJobApplication(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	uuid, err := common.ToUUID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var body models.CreateJobApplicationReqBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	jobApplication, err := h.queries.CreateJobApplication(h.ctx, db.CreateJobApplicationParams{
		UserID:        uuid,
		CompanyName:   body.CompanyName,
		JobTitle:      body.JobTitle,
		DateApplied:   pgtype.Timestamp{Time: body.DateApplied, Valid: true},
		Status:        db.Status(body.Status),
		MinSalary:     pgtype.Float8{Float64: body.MinSalary, Valid: true},
		MaxSalary:     pgtype.Float8{Float64: body.MaxSalary, Valid: true},
		JobPostingUrl: pgtype.Text{String: body.JobPostingURL, Valid: true},
		Notes:         pgtype.Text{String: body.Notes, Valid: true},
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	createJobApplicationResBody := models.NewCreateJobApplicationResBody(jobApplication)

	c.JSON(http.StatusCreated, createJobApplicationResBody)
}

// TODO
func (h *Handler) UpdateJobApplication(c *gin.Context) {}

// TODO
func (h *Handler) DeleteJobApplication(c *gin.Context) {}
