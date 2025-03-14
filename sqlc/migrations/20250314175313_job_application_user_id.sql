-- +goose Up
-- +goose StatementBegin
ALTER TABLE job_applications ADD COLUMN user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE job_applications DROP COLUMN user_id;
-- +goose StatementEnd
