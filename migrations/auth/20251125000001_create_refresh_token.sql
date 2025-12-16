-- +goose Up
CREATE TABLE refresh_tokens (
	id UUID PRIMARY KEY,
	user_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
	token_hash VARCHAR(255) UNIQUE NOT NULL,
	expires_at TIMESTAMP NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;
DROP TABLE IF EXISTS refresh_tokens;
