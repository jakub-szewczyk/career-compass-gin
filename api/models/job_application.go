package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
)

type Sort string

const (
	CompanyNameAsc  Sort = "company_name"
	CompanyNameDesc Sort = "-company_name"
	JobTitleAsc     Sort = "job_title"
	JobTitleDesc    Sort = "-job_title"
	DateAppliedAsc  Sort = "date_applied"
	DateAppliedDesc Sort = "-date_applied"
	StatusAsc       Sort = "status"
	StatusDesc      Sort = "-status"
	SalaryAsc       Sort = "salary"
	SalaryDesc      Sort = "-salary"
	IsRepliedAsc    Sort = "is_replied"
	IsRepliedDesc   Sort = "-is_replied"
)

type JobApplicationsQueryParams struct {
	Page                  int       `form:"page" binding:"min=0"`
	Size                  int       `form:"size" binding:"min=0"`
	Sort                  Sort      `form:"sort" binding:"omitempty,oneof=company_name -company_name job_title -job_title date_applied -date_applied status -status salary -salary is_replied -is_replied"`
	CompanyNameOrJobTitle string    `form:"company_name_or_job_title" binding:"omitempty"`
	DateApplied           string    `form:"date_applied" binding:"omitempty,datetime=2006-01-02"`
	Status                db.Status `form:"status" binding:"omitempty,oneof=IN_PROGRESS REJECTED ACCEPTED"`
}

type jobApplicationEntry struct {
	ID            string    `json:"id" example:"f4d15edc-e780-42b5-957d-c4352401d9ca"`
	CompanyName   string    `json:"companyName" example:"Evil Corp Inc."`
	JobTitle      string    `json:"jobTitle" example:"Software Engineer"`
	DateApplied   time.Time `json:"dateApplied" example:"2025-03-14T12:34:56Z"`
	Status        db.Status `json:"status" example:"IN_PROGRESS"`
	IsReplied     bool      `json:"isReplied" example:"false"`
	MinSalary     float64   `json:"minSalary,omitempty" example:"50000.00"`
	MaxSalary     float64   `json:"maxSalary,omitempty" example:"70000.00"`
	JobPostingURL string    `json:"jobPostingURL,omitempty" example:"https://glassbore.com/jobs/swe420692137"`
}

type JobApplicationsResBody struct {
	Page  int                   `json:"page" example:"0"`
	Size  int                   `json:"size" example:"10"`
	Total int                   `json:"total" example:"100"`
	Data  []jobApplicationEntry `json:"data"`
}

func NewJobApplicationsResBody(page, size int, jobApplications []db.GetJobApplicationsRow) JobApplicationsResBody {
	data := []jobApplicationEntry{}

	for _, jobApplication := range jobApplications {
		data = append(data, jobApplicationEntry{
			ID:            jobApplication.ID.String(),
			CompanyName:   jobApplication.CompanyName,
			JobTitle:      jobApplication.JobTitle,
			DateApplied:   jobApplication.DateApplied.Time.UTC(),
			Status:        jobApplication.Status,
			IsReplied:     jobApplication.IsReplied,
			MinSalary:     jobApplication.MinSalary.Float64,
			MaxSalary:     jobApplication.MaxSalary.Float64,
			JobPostingURL: jobApplication.JobPostingUrl.String,
		})
	}

	total := 0

	if len(jobApplications) > 0 {
		total = int(jobApplications[0].Total)
	}

	return JobApplicationsResBody{
		Page:  page,
		Size:  size,
		Total: total,
		Data:  data,
	}
}

type JobApplicationResBody struct {
	ID            string    `json:"id" example:"f4d15edc-e780-42b5-957d-c4352401d9ca"`
	CompanyName   string    `json:"companyName" example:"Evil Corp Inc."`
	JobTitle      string    `json:"jobTitle" example:"Software Engineer"`
	DateApplied   time.Time `json:"dateApplied" example:"2025-03-14T12:34:56Z"`
	Status        db.Status `json:"status" example:"IN_PROGRESS"`
	IsReplied     bool      `json:"isReplied" example:"false"`
	MinSalary     float64   `json:"minSalary,omitempty" example:"50000.00"`
	MaxSalary     float64   `json:"maxSalary,omitempty" example:"70000.00"`
	JobPostingURL string    `json:"jobPostingURL,omitempty" example:"https://glassbore.com/jobs/swe420692137"`
	Notes         string    `json:"notes,omitempty" example:"Follow up in two weeks"`
}

func NewJobApplicationResBody(jobApplication db.GetJobApplicationRow) JobApplicationResBody {
	return JobApplicationResBody{
		ID:            jobApplication.ID.String(),
		CompanyName:   jobApplication.CompanyName,
		JobTitle:      jobApplication.JobTitle,
		DateApplied:   jobApplication.DateApplied.Time.UTC(),
		Status:        jobApplication.Status,
		IsReplied:     jobApplication.IsReplied,
		MinSalary:     jobApplication.MinSalary.Float64,
		MaxSalary:     jobApplication.MaxSalary.Float64,
		JobPostingURL: jobApplication.JobPostingUrl.String,
		Notes:         jobApplication.Notes.String,
	}
}

type CreateJobApplicationReqBody struct {
	CompanyName   string    `json:"companyName" binding:"required" example:"Evil Corp Inc."`
	JobTitle      string    `json:"jobTitle" binding:"required" example:"Software Engineer"`
	DateApplied   time.Time `json:"dateApplied" binding:"required" example:"2025-03-14T12:34:56Z"`
	Status        db.Status `json:"status" binding:"required,oneof=IN_PROGRESS REJECTED ACCEPTED" example:"IN_PROGRESS"`
	MinSalary     float64   `json:"minSalary,omitempty" binding:"omitempty,gte=0" example:"50000.00"`
	MaxSalary     float64   `json:"maxSalary,omitempty" binding:"omitempty,gte=0" example:"70000.00"`
	JobPostingURL string    `json:"jobPostingURL,omitempty" example:"https://glassbore.com/jobs/swe420692137"`
	Notes         string    `json:"notes,omitempty" example:"Follow up in two weeks"`
}

func NewCreateJobApplicationReqBody(companyName, jobTitle string, dateApplied time.Time, status db.Status, minSalary, maxSalary float64, jobPostingURL, notes string) CreateJobApplicationReqBody {
	return CreateJobApplicationReqBody{
		CompanyName:   companyName,
		JobTitle:      jobTitle,
		DateApplied:   dateApplied,
		Status:        status,
		MinSalary:     minSalary,
		MaxSalary:     maxSalary,
		JobPostingURL: jobPostingURL,
		Notes:         notes,
	}
}

type CreateJobApplicationResBody struct {
	ID            string    `json:"id" example:"f4d15edc-e780-42b5-957d-c4352401d9ca"`
	CompanyName   string    `json:"companyName" example:"Evil Corp Inc."`
	JobTitle      string    `json:"jobTitle" example:"Software Engineer"`
	DateApplied   time.Time `json:"dateApplied" example:"2025-03-14T12:34:56Z"`
	Status        db.Status `json:"status" example:"IN_PROGRESS"`
	IsReplied     bool      `json:"isReplied" example:"false"`
	MinSalary     float64   `json:"minSalary,omitempty" example:"50000.00"`
	MaxSalary     float64   `json:"maxSalary,omitempty" example:"70000.00"`
	JobPostingURL string    `json:"jobPostingURL,omitempty" example:"https://glassbore.com/jobs/swe420692137"`
	Notes         string    `json:"notes,omitempty" example:"Follow up in two weeks"`
}

func NewCreateJobApplicationResBody(jobApplication db.CreateJobApplicationRow) CreateJobApplicationResBody {
	return CreateJobApplicationResBody{
		ID:            jobApplication.ID.String(),
		CompanyName:   jobApplication.CompanyName,
		JobTitle:      jobApplication.JobTitle,
		DateApplied:   jobApplication.DateApplied.Time.UTC(),
		Status:        jobApplication.Status,
		IsReplied:     jobApplication.IsReplied,
		MinSalary:     jobApplication.MinSalary.Float64,
		MaxSalary:     jobApplication.MaxSalary.Float64,
		JobPostingURL: jobApplication.JobPostingUrl.String,
		Notes:         jobApplication.Notes.String,
	}
}

type UpdateJobApplicationReqBody struct {
	CompanyName   string     `json:"companyName,omitempty" example:"Evil Corp Inc."`
	JobTitle      string     `json:"jobTitle,omitempty" example:"Software Engineer"`
	DateApplied   *time.Time `json:"dateApplied,omitempty" example:"2025-03-14T12:34:56Z"`
	Status        *db.Status `json:"status,omitempty" binding:"omitempty,oneof=IN_PROGRESS REJECTED ACCEPTED" example:"IN_PROGRESS"`
	IsReplied     *bool      `json:"isReplied,omitempty" example:"false"`
	MinSalary     *float64   `json:"minSalary,omitempty" binding:"omitempty,gte=0" example:"50000.00"`
	MaxSalary     *float64   `json:"maxSalary,omitempty" binding:"omitempty,gte=0" example:"70000.00"`
	JobPostingURL string     `json:"jobPostingURL,omitempty" example:"https://glassbore.com/jobs/swe420692137"`
	Notes         string     `json:"notes,omitempty" example:"Follow up in two weeks"`
}

func NewUpdateJobApplicationReqBody(companyName, jobTitle string, dateApplied *time.Time, status *db.Status, isReplied *bool, minSalary, maxSalary *float64, jobPostingURL, notes string) UpdateJobApplicationReqBody {
	return UpdateJobApplicationReqBody{
		CompanyName:   companyName,
		JobTitle:      jobTitle,
		DateApplied:   dateApplied,
		Status:        status,
		IsReplied:     isReplied,
		MinSalary:     minSalary,
		MaxSalary:     maxSalary,
		JobPostingURL: jobPostingURL,
		Notes:         notes,
	}
}

type UpdateJobApplicationResBody struct {
	ID            string    `json:"id" example:"f4d15edc-e780-42b5-957d-c4352401d9ca"`
	CompanyName   string    `json:"companyName" example:"Evil Corp Inc."`
	JobTitle      string    `json:"jobTitle" example:"Software Engineer"`
	DateApplied   time.Time `json:"dateApplied" example:"2025-03-14T12:34:56Z"`
	Status        db.Status `json:"status" example:"IN_PROGRESS"`
	IsReplied     bool      `json:"isReplied" example:"false"`
	MinSalary     float64   `json:"minSalary,omitempty" example:"50000.00"`
	MaxSalary     float64   `json:"maxSalary,omitempty" example:"70000.00"`
	JobPostingURL string    `json:"jobPostingURL,omitempty" example:"https://glassbore.com/jobs/swe420692137"`
	Notes         string    `json:"notes,omitempty" example:"Follow up in two weeks"`
}

func NewUpdateJobApplicationResBody(jobApplication db.UpdateJobApplicationRow) UpdateJobApplicationResBody {
	return UpdateJobApplicationResBody{
		ID:            jobApplication.ID.String(),
		CompanyName:   jobApplication.CompanyName,
		JobTitle:      jobApplication.JobTitle,
		DateApplied:   jobApplication.DateApplied.Time.UTC(),
		Status:        jobApplication.Status,
		IsReplied:     jobApplication.IsReplied,
		MinSalary:     jobApplication.MinSalary.Float64,
		MaxSalary:     jobApplication.MaxSalary.Float64,
		JobPostingURL: jobApplication.JobPostingUrl.String,
		Notes:         jobApplication.Notes.String,
	}
}

func NewUpdateJobApplicationParams(jobApplicationId, userId pgtype.UUID, body UpdateJobApplicationReqBody) db.UpdateJobApplicationParams {
	params := db.UpdateJobApplicationParams{
		ID:            jobApplicationId,
		UserID:        userId,
		CompanyName:   pgtype.Text{String: body.CompanyName, Valid: true},
		JobTitle:      pgtype.Text{String: body.JobTitle, Valid: true},
		JobPostingUrl: pgtype.Text{String: body.JobPostingURL, Valid: true},
		Notes:         pgtype.Text{String: body.Notes, Valid: true},
	}

	if body.DateApplied != nil {
		params.DateApplied = pgtype.Timestamptz{Time: *body.DateApplied, Valid: true}
	}
	if body.Status != nil {
		params.Status = db.NullStatus{Status: *body.Status, Valid: true}
	}
	if body.IsReplied != nil {
		params.IsReplied = pgtype.Bool{Bool: *body.IsReplied, Valid: true}
	}
	if body.MinSalary != nil {
		params.MinSalary = pgtype.Float8{Float64: *body.MinSalary, Valid: true}
	}
	if body.MaxSalary != nil {
		params.MaxSalary = pgtype.Float8{Float64: *body.MaxSalary, Valid: true}
	}

	return params
}
