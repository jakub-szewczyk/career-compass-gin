// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	Email     string      `json:"email"`
	Password  string      `json:"password"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	ID        pgtype.UUID `json:"id"`
}
