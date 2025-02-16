package models

import (
	"errors"

	"github.com/jakub-szewczyk/career-compass-gin/sqlc/db"
)

type ProfileResBody struct {
	ID                string `json:"id" example:"f4d15edc-e780-42b5-957d-c4352401d9ca"`
	FirstName         string `json:"firstName" example:"John"`
	LastName          string `json:"lastName" example:"Doe"`
	Email             string `json:"email" example:"john.doe@example.com"`
	IsEmailVerified   bool   `json:"isEmailVerified" example:"true"`
	VerificationToken string `json:"verificationToken" example:"2cc313c8b72f8e5b725e07130d0b851811f2e60c8b19f085b3aa58d1516ef767"`
}

type AnyUser interface{}

func NewProfileResBody(user AnyUser) (*ProfileResBody, error) {
	switch user := user.(type) {
	case db.CreateUserRow:
		return &ProfileResBody{
			ID:                user.ID.String(),
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Email:             user.Email,
			IsEmailVerified:   user.IsEmailVerified.Bool,
			VerificationToken: user.VerificationToken,
		}, nil
	case db.GetUserOnSignInRow:
		return &ProfileResBody{
			ID:                user.ID.String(),
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Email:             user.Email,
			IsEmailVerified:   user.IsEmailVerified.Bool,
			VerificationToken: user.VerificationToken,
		}, nil
	case db.GetUserByIdRow:
		return &ProfileResBody{
			ID:                user.ID.String(),
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Email:             user.Email,
			IsEmailVerified:   user.IsEmailVerified.Bool,
			VerificationToken: user.VerificationToken,
		}, nil
	}
	return nil, errors.New("unhandled user type")
}
