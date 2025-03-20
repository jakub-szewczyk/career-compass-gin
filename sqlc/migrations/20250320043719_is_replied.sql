-- +goose Up
-- +goose StatementBegin
ALTER TABLE job_applications ADD COLUMN is_replied BOOLEAN NOT NULL DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE job_applications DROP COLUMN is_replied;
-- +goose StatementEnd
