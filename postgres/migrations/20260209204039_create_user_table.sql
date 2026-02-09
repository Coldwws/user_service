-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    first_name TEXT,
    last_name TEXT,
    password TEXT,
    phone_number TEXT,
    email TEXT,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd