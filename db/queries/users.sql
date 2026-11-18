-- name: RegisterUser :one
INSERT INTO users (
    email,
    pass,
    username,
    created_at,
    updated_at
) VALUES (
    sqlc.arg(email),
    sqlc.arg(pass),
    sqlc.arg(username),
    DATETIME('now'),
    DATETIME('now')
) RETURNING *;