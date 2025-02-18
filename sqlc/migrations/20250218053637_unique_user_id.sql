-- +goose Up
ALTER TABLE verification_tokens ADD CONSTRAINT unique_user_id UNIQUE (user_id);

-- +goose Down
ALTER TABLE verification_tokens DROP CONSTRAINT unique_user_id;
