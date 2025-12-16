-- +goose Up
CREATE TABLE auth_users (
	id UUID PRIMARY KEY,
	password_hash VARCHAR(255) NOT NULL,
	email VARCHAR(100) UNIQUE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS auth_users;