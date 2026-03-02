-- +goose Up
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    first_name TEXT not null,
    last_name TEXT not null,
    password TEXT not null,
    phone_number TEXT not null,
    email TEXT not null,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd