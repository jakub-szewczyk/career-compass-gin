package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	common "github.com/jakub-szewczyk/career-compass-gin/utils"
)

// JobApplications godoc
//
//	@Summary		Get job applications
//	@Description	Retrieves a list of job applications with support for sorting, filtering, and pagination
//
//	@Security		BearerAuth
//
//	@Tags			Job application
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page number (zero-indexed)"	minimum(0)	default(0)
//	@Param			size	query		int	false	"Page size"						minimum(0)	default(10)
//	@Failure		400		{object}	models.Error
//	@Failure		500		{object}	models.Error
//	@Success		200		{object}	models.JobApplicationsResBody
//	@Router			/job-applications [get]
func (h *Handler) JobApplications(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	uuid, err := common.ToUUID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var queryParams models.JobApplicationsQueryParams

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if c.Query("size") == "" {
		queryParams.Size = 10
	}

	// TODO:
	// - Support sorting by company name, job title, date applied, status, salary, and replied column.
	// - Support filtering by company name, job title, date applied, and status.
	jobApplications, err := h.queries.GetJobApplications(h.ctx, db.GetJobApplicationsParams{
		Limit:  int32(queryParams.Size),
		Offset: int32(queryParams.Page * queryParams.Size),
		UserID: uuid,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resBody := models.NewJobApplicationsResBody(queryParams.Page, queryParams.Size, jobApplications)

	c.JSON(http.StatusOK, resBody)
}

// JobApplication godoc
//
//	@Summary		Retrieve job application details
//	@Description	Fetches the details of a specific job application by its id
//
//	@Security		BearerAuth
//
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

	resBody := models.NewJobApplicationResBody(jobApplication)

	c.JSON(http.StatusOK, resBody)
}

// CreateJobApplication godoc
//
//	@Summary		Submit a new job application
//	@Description	Processes and creates a new job application with the provided data
//
//	@Security		BearerAuth
//
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
		DateApplied:   pgtype.Timestamptz{Time: body.DateApplied, Valid: true},
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

	resBody := models.NewCreateJobApplicationResBody(jobApplication)

	c.JSON(http.StatusCreated, resBody)
}

// TODO
func (h *Handler) UpdateJobApplication(c *gin.Context) {}

// TODO
func (h *Handler) DeleteJobApplication(c *gin.Context) {}
