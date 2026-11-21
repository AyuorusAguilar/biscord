-- +goose Up
CREATE TABLE sessions (
	id INTEGER PRIMARY KEY,
	userID INTEGER NOT NULL,
	refresh_token_hash BLOB NOT NULL,
	expires_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL,
	manual_expiration BOOLEAN NOT NULL DEFAULT FALSE,
	CONSTRAINT FK_userID FOREIGN KEY (userID) REFERENCES users(id) ON DELETE CASCADEz
);
-- +goose Down
DROP TABLE sessions;