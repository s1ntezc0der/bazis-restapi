-- +goose Up
ALTER TABLE tasks ADD COLUMN priority INT DEFAULT 0;

-- +goose Down
ALTER TABLE tasks DROP COLUMN priority;