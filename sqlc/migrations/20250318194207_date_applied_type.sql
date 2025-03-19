-- +goose Up
-- +goose StatementBegin
ALTER TABLE job_applications ALTER COLUMN date_applied TYPE TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE job_applications ALTER COLUMN date_applied TYPE TIMESTAMP;
-- +goose StatementEnd
