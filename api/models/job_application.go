package models

import (
	"time"

	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
)

type JobApplicationResBody struct {
	ID            string    `json:"id" example:"f4d15edc-e780-42b5-957d-c4352401d9ca"`
	CompanyName   string    `json:"companyName" binding:"required" example:"Evil Corp Inc."`
	JobTitle      string    `json:"jobTitle" binding:"required" example:"Software Engineer"`
	DateApplied   time.Time `json:"dateApplied" binding:"required" example:"2025-03-14T12:34:56Z"`
	Status        db.Status `json:"status" binding:"required" example:"IN_PROGRESS"`
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
	Status        db.Status `json:"status" binding:"required" example:"IN_PROGRESS"`
	MinSalary     float64   `json:"minSalary,omitempty" example:"50000.00"`
	MaxSalary     float64   `json:"maxSalary,omitempty" example:"70000.00"`
	JobPostingURL string    `json:"jobPostingURL,omitempty" example:"https://glassbore.com/jobs/swe420692137"`
	Notes         string    `json:"notes,omitempty" example:"Follow up in two weeks"`
}

func NewCreateJobApplicationReqBody(companyName, jobTitle string, dateApplied time.Time, status db.Status, minSalary, maxSalary float64, jobPostingURL string, notes string) CreateJobApplicationReqBody {
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
		MinSalary:     jobApplication.MinSalary.Float64,
		MaxSalary:     jobApplication.MaxSalary.Float64,
		JobPostingURL: jobApplication.JobPostingUrl.String,
		Notes:         jobApplication.Notes.String,
	}
}
