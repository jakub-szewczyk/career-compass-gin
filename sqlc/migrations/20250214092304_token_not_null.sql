-- +goose Up
ALTER TABLE verification_tokens ALTER COLUMN token SET NOT NULL;
