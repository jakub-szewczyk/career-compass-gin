package models

import (
	"time"

	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
)

const (
	TitleAsc      = "title"
	TitleDesc     = "-title"
	CreatedAtAsc  = "created_at"
	CreatedAtDesc = "-created_at"
	UpdatedAtAsc  = "updated_at"
	UpdatedAtDesc = "-updated_at"
)

type ResumesQueryParams struct {
	Page  int    `form:"page" binding:"min=0"`
	Size  int    `form:"size" binding:"min=0"`
	Sort  Sort   `form:"sort" binding:"omitempty,oneof=title -title created_at -created_at updated_at -updated_at"`
	Title string `form:"title" binding:"omitempty"`
}

type resumeEntry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ResumesResBody struct {
	Page  int           `json:"page" example:"0"`
	Size  int           `json:"size" example:"10"`
	Total int           `json:"total" example:"100"`
	Data  []resumeEntry `json:"data"`
}

func NewResumesResBody(page, size int, resumes []db.GetResumesRow) ResumesResBody {
	data := []resumeEntry{}

	for _, resume := range resumes {
		data = append(data, resumeEntry{
			ID:        resume.ID.String(),
			Title:     resume.Title,
			CreatedAt: resume.CreatedAt.Time.UTC(),
			UpdatedAt: resume.UpdatedAt.Time.UTC(),
		})
	}

	total := 0

	if len(resumes) > 0 {
		total = int(resumes[0].Total)
	}

	return ResumesResBody{
		Page:  page,
		Size:  size,
		Total: total,
		Data:  data,
	}
}

type CreateResumeReqBody struct {
	Title string `json:"title,omitempty" example:"Evil Corp Inc. personalized"`
}

func NewCreateResumeReqBody(title string) CreateResumeReqBody {
	return CreateResumeReqBody{
		Title: title,
	}
}

type CreateResumeResBody struct {
	ID    string `json:"id" example:"f4d15edc-e780-42b5-957d-c4352401d9ca"`
	Title string `json:"title,omitempty" example:"Evil Corp Inc. personalized"`
}

func NewCreateResumeResBody(resume db.CreateResumeRow) CreateResumeResBody {
	return CreateResumeResBody{ID: resume.ID.String(), Title: resume.Title}
}

type DeleteResumeResBody struct {
	ID    string `json:"id" example:"f4d15edc-e780-42b5-957d-c4352401d9ca"`
	Title string `json:"title,omitempty" example:"Evil Corp Inc. personalized"`
}

func NewDeleteResumeResBody(resume db.DeleteResumeRow) DeleteResumeResBody {
	return DeleteResumeResBody{ID: resume.ID.String(), Title: resume.Title}
}
