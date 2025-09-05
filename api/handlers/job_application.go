package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jakub-szewczyk/career-compass-gin/api/models"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
	"github.com/jakub-szewczyk/career-compass-gin/utils"
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
//	@Param			page						query		int		false	"Page number (zero-indexed)"	minimum(0)																																			default(0)
//	@Param			size						query		int		false	"Page size"						minimum(0)																																			default(10)
//	@Param			sort						query		string	false	"Sortable column name"			Enums(company_name, -company_name, job_title, -job_title, date_applied, -date_applied, status, -status, salary, -salary, is_replied, -is_replied)	default(-date_applied)
//	@Param			company_name_or_job_title	query		string	false	"Company name or job title"
//	@Param			date_applied				query		string	false	"Date applied"
//	@Param			status						query		string	false	"Status"	Enums(IN_PROGRESS, REJECTED, ACCEPTED)
//	@Failure		400							{object}	models.Error
//	@Failure		500							{object}	models.Error
//	@Success		200							{object}	models.JobApplicationsResBody
//	@Router			/job-applications [get]
func (h *Handler) JobApplications(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	uuid, err := utils.ToUUID(userId)
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

	if c.Query("sort") == "" {
		queryParams.Sort = models.DateAppliedDesc
	}

	var dateApplied time.Time
	if queryParams.DateApplied != "" {
		dateApplied, err = time.Parse(time.DateOnly, queryParams.DateApplied)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	jobApplications, err := h.queries.GetJobApplications(h.ctx, db.GetJobApplicationsParams{
		UserID: uuid,

		Limit:  int32(queryParams.Size),
		Offset: int32(queryParams.Page * queryParams.Size),

		CompanyNameAsc:  queryParams.Sort == models.CompanyNameAsc,
		CompanyNameDesc: queryParams.Sort == models.CompanyNameDesc,
		JobTitleAsc:     queryParams.Sort == models.JobTitleAsc,
		JobTitleDesc:    queryParams.Sort == models.JobTitleDesc,
		DateAppliedAsc:  queryParams.Sort == models.DateAppliedAsc,
		DateAppliedDesc: queryParams.Sort == models.DateAppliedDesc,
		StatusAsc:       queryParams.Sort == models.StatusAsc,
		StatusDesc:      queryParams.Sort == models.StatusDesc,
		SalaryAsc:       queryParams.Sort == models.SalaryAsc,
		SalaryDesc:      queryParams.Sort == models.SalaryDesc,
		IsRepliedAsc:    queryParams.Sort == models.IsRepliedAsc,
		IsRepliedDesc:   queryParams.Sort == models.IsRepliedDesc,

		CompanyNameOrJobTitle: queryParams.CompanyNameOrJobTitle,
		DateApplied:           utils.NullifyTime(dateApplied),
		Status:                queryParams.Status,
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

	uuid, err := utils.ToUUID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	jobApplicationId, err := utils.ToUUID(c.Param("jobApplicationId"))
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

	uuid, err := utils.ToUUID(userId)
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

// UpdateJobApplication godoc
//
//	@Summary		Update a job application
//	@Description	Updates an existing job application with the provided details
//
//	@Security		BearerAuth
//
//	@Tags			Job application
//	@Accept			json
//	@Produce		json
//	@Param			jobApplicationId	path		string								true	"Job application uuid"
//	@Param			body				body		models.UpdateJobApplicationReqBody	true	"Job application details"
//	@Failure		400					{object}	models.Error
//	@Failure		500					{object}	models.Error
//	@Success		200					{object}	models.UpdateJobApplicationResBody
//	@Router			/job-applications/{jobApplicationId} [put]
func (h *Handler) UpdateJobApplication(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	uuid, err := utils.ToUUID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	jobApplicationId, err := utils.ToUUID(c.Param("jobApplicationId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var body models.UpdateJobApplicationReqBody

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	params := models.NewUpdateJobApplicationParams(jobApplicationId, uuid, body)

	jobApplication, err := h.queries.UpdateJobApplication(h.ctx, params)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resBody := models.NewUpdateJobApplicationResBody(jobApplication)

	c.JSON(http.StatusOK, resBody)
}

// DeleteJobApplication godoc
//
//	@Summary		Delete a job application
//	@Description	Deletes an existing job application
//
//	@Security		BearerAuth
//
//	@Tags			Job application
//	@Accept			json
//	@Produce		json
//	@Param			jobApplicationId	path		string	true	"Job application uuid"
//	@Failure		500					{object}	models.Error
//	@Success		200					{object}	models.DeleteJobApplicationResBody
//	@Router			/job-applications/{jobApplicationId} [delete]
func (h *Handler) DeleteJobApplication(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	uuid, err := utils.ToUUID(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	jobApplicationId, err := utils.ToUUID(c.Param("jobApplicationId"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	jobApplication, err := h.queries.DeleteJobApplication(h.ctx, db.DeleteJobApplicationParams{
		ID:     jobApplicationId,
		UserID: uuid,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resBody := models.NewDeleteJobApplicationResBody(jobApplication)

	c.JSON(http.StatusOK, resBody)
}
