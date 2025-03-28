package utils

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToUUID(s string) (pgtype.UUID, error) {
	var uuid pgtype.UUID
	err := uuid.Scan(s)
	return uuid, err
}

func NullifyTime(date time.Time) interface{} {
	if date.IsZero() {
		return nil
	}
	return date
}
