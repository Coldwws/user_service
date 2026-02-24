-- +goose Up
ALTER TABLE users
ALTER COLUMN created_at
SET DEFAULT now(),
ALTER COLUMN updated_at
SET DEFAULT now();

-- +goose Down
ALTER TABLE users
ALTER COLUMN created_at
DROP DEFAULT,
ALTER COLUMN updated_at
DROP DEFAULT;