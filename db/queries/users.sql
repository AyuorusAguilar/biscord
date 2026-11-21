-- name: RegisterUser :one
INSERT INTO users (
	username,
	pass,
	created_at,
	updated_at
) VALUES (
	sqlc.arg(username),
	sqlc.arg(pass),
	DATETIME('now'),
	DATETIME('now')
) RETURNING username, created_at;

-- name: UsernameExists :one
SELECT EXISTS (
	SELECT 1 FROM users
	WHERE username = sqlc.arg(username)
) as existe; /* Probably all queries should be usable in the sqlite3 terminal, that's the reason of this alias */

-- name: GetIdByUsername :one
SELECT id FROM users WHERE username = sqlc.arg(username);