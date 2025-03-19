-- +goose Up
ALTER TABLE users ALTER COLUMN created_at TYPE timestamptz;
ALTER TABLE users ALTER COLUMN updated_at TYPE timestamptz;
ALTER TABLE verification_tokens ALTER COLUMN expires_at TYPE timestamptz;
ALTER TABLE verification_tokens ALTER COLUMN created_at TYPE timestamptz;
ALTER TABLE verification_tokens ALTER COLUMN updated_at TYPE timestamptz;
ALTER TABLE password_reset_tokens ALTER COLUMN expires_at TYPE timestamptz;
ALTER TABLE password_reset_tokens ALTER COLUMN created_at TYPE timestamptz;
ALTER TABLE password_reset_tokens ALTER COLUMN updated_at TYPE timestamptz;
ALTER TABLE job_applications ALTER COLUMN created_at TYPE timestamptz;
ALTER TABLE job_applications ALTER COLUMN updated_at TYPE timestamptz;

-- +goose Down
ALTER TABLE users ALTER COLUMN created_at TYPE timestamp;
ALTER TABLE users ALTER COLUMN updated_at TYPE timestamp;
ALTER TABLE verification_tokens ALTER COLUMN expires_at TYPE timestamp;
ALTER TABLE verification_tokens ALTER COLUMN created_at TYPE timestamp;
ALTER TABLE verification_tokens ALTER COLUMN updated_at TYPE timestamp;
ALTER TABLE password_reset_tokens ALTER COLUMN expires_at TYPE timestamp;
ALTER TABLE password_reset_tokens ALTER COLUMN created_at TYPE timestamp;
ALTER TABLE password_reset_tokens ALTER COLUMN updated_at TYPE timestamp;
ALTER TABLE job_applications ALTER COLUMN created_at TYPE timestamp;
ALTER TABLE job_applications ALTER COLUMN updated_at TYPE timestamp;
