-- +goose Up
-- +goose StatementBegin
ALTER TABLE job_applications DROP CONSTRAINT IF EXISTS job_applications_user_id_key;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE job_applications ADD CONSTRAINT job_applications_user_id_key UNIQUE (user_id);
-- +goose StatementEnd
