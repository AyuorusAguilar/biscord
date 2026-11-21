-- name: RegisterSession :one
INSERT INTO sessions (
	userID,
	refresh_token_hash,
	expires_at,
	created_at
) values (
	sqlc.arg(userID),
	sqlc.arg(refresh_token_hash),
	sqlc.arg(expires_at),
	DATETIME('now')
) RETURNING *;

-- name: GetLastSession :one
SELECT * FROM sessions 
	WHERE 
		userID = sqlc.arg(userID)
	ORDER BY created_at DESC
	LIMIT 1;

-- name: ExpireSession :one
UPDATE sessions
	SET manual_expiration = TRUE
	WHERE id = sqlc.arg(sessionID)
	RETURNING *;
-- name: GetSessionByRefreshToken :one
SELECT * FROM sessions
	WHERE refresh_token_hash = sqlc.arg(refresh_token_hash)
	LIMIT 1;