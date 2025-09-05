package models

import "github.com/jakub-szewczyk/career-compass-gin/sqlc/db"

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
