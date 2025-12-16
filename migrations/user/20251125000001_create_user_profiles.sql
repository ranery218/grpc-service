-- +goose Up
CREATE TABLE user_profiles (
	user_id UUID PRIMARY KEY NOT NULL,
	username VARCHAR(50) UNIQUE NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS user_profiles;