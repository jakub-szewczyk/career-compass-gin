-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_updated_at_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TYPE status AS ENUM ('IN_PROGRESS', 'REJECTED', 'ACCEPTED');
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE job_applications (
  id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  company_name    TEXT NOT NULL,
  job_title       TEXT NOT NULL,
  date_applied    TIMESTAMP NOT NULL,
  status          status NOT NULL DEFAULT 'IN_PROGRESS',
  min_salary      DOUBLE PRECISION,
  max_salary      DOUBLE PRECISION,
  job_posting_url TEXT,
  notes           TEXT,
  created_at      TIMESTAMP DEFAULT NOW(),
  updated_at      TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER set_job_application_updated_at_timestamp
BEFORE UPDATE ON job_applications
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS status;
DROP TABLE IF EXISTS job_applications;
DROP FUNCTION IF EXISTS set_updated_at_timestamp;
-- +goose StatementEnd
