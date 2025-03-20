// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type Status string

const (
	StatusINPROGRESS Status = "IN_PROGRESS"
	StatusREJECTED   Status = "REJECTED"
	StatusACCEPTED   Status = "ACCEPTED"
)

func (e *Status) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = Status(s)
	case string:
		*e = Status(s)
	default:
		return fmt.Errorf("unsupported scan type for Status: %T", src)
	}
	return nil
}

type NullStatus struct {
	Status Status `json:"status"`
	Valid  bool   `json:"valid"` // Valid is true if Status is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullStatus) Scan(value interface{}) error {
	if value == nil {
		ns.Status, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.Status.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.Status), nil
}

type JobApplication struct {
	ID            pgtype.UUID        `json:"id"`
	CompanyName   string             `json:"companyName"`
	JobTitle      string             `json:"jobTitle"`
	DateApplied   pgtype.Timestamptz `json:"dateApplied"`
	Status        Status             `json:"status"`
	MinSalary     pgtype.Float8      `json:"minSalary"`
	MaxSalary     pgtype.Float8      `json:"maxSalary"`
	JobPostingUrl pgtype.Text        `json:"jobPostingUrl"`
	Notes         pgtype.Text        `json:"notes"`
	CreatedAt     pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt     pgtype.Timestamptz `json:"updatedAt"`
	UserID        pgtype.UUID        `json:"userId"`
	IsReplied     bool               `json:"isReplied"`
}

type PasswordResetToken struct {
	ID        pgtype.UUID        `json:"id"`
	UserID    pgtype.UUID        `json:"userId"`
	Token     string             `json:"token"`
	ExpiresAt pgtype.Timestamptz `json:"expiresAt"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
}

type User struct {
	Email           string             `json:"email"`
	Password        string             `json:"password"`
	FirstName       string             `json:"firstName"`
	LastName        string             `json:"lastName"`
	ID              pgtype.UUID        `json:"id"`
	CreatedAt       pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt       pgtype.Timestamptz `json:"updatedAt"`
	IsEmailVerified pgtype.Bool        `json:"isEmailVerified"`
}

type VerificationToken struct {
	ID        pgtype.UUID        `json:"id"`
	UserID    pgtype.UUID        `json:"userId"`
	Token     string             `json:"token"`
	ExpiresAt pgtype.Timestamptz `json:"expiresAt"`
	CreatedAt pgtype.Timestamptz `json:"createdAt"`
	UpdatedAt pgtype.Timestamptz `json:"updatedAt"`
}
