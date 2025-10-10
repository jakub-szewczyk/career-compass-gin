-- +goose Up
CREATE TABLE resumes (
  id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id         UUID REFERENCES users(id) ON DELETE CASCADE,
  title           TEXT NOT NULL UNIQUE,
  created_at      TIMESTAMPTZ DEFAULT NOW(),
  updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TRIGGER set_resume_updated_at_timestamp
BEFORE UPDATE ON resumes
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_timestamp();

-- +goose Down
DROP TABLE IF EXISTS resumes;
DROP TRIGGER IF EXISTS set_resume_updated_at_timestamp ON resumes;