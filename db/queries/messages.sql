-- name: CreateDMChannel :one 
INSERT INTO channels (
    type, 
    created_at
) VALUES(
    "DM",
    DATETIME('now')
) RETURNING id;

-- name: SendMessage :one
INSERT INTO messages (
    senderID,
    channelID,
    content,
    created_at
) VALUES (
    sqlc.arg(senderID),
    sqlc.arg(channelID),
    sqlc.arg(content),
    DATETIME('now')
) RETURNING *;

-- name: GetMessagesFromChannelByUserID :many
SELECT 
    id,
    content,
    created_at
FROM messages
WHERE 
    senderID = sqlc.arg(senderID) AND
    channelID = sqlc.arg(channelID);