-- +goose Up
ALTER TABLE resumes DROP CONSTRAINT resumes_title_key;
ALTER TABLE resumes ADD CONSTRAINT resumes_title_user_id_key UNIQUE (title, user_id);

-- +goose Down
ALTER TABLE resumes DROP CONSTRAINT resumes_title_user_id_key;
ALTER TABLE resumes ADD CONSTRAINT resumes_title_key UNIQUE (title);
