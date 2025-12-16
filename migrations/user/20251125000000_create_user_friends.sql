-- +goose Up
CREATE TABLE user_friends (
	user_id UUID PRIMARY KEY NOT NULL,
	friend_ids UUID[]
);

-- +goose Down
DROP TABLE IF EXISTS user_friends;