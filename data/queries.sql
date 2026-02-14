-- db/queries.sql

-- name: CreateUser :execresult
INSERT INTO users (username, password)
VALUES (?, ?);

-- name: GetUser :one
SELECT id, username, password
FROM users
WHERE username = ?;

-- name: GetUserByID :one
SELECT id, username, password
FROM users
WHERE id = ?;

-- name: GetAllUsers :many
SELECT id, username
FROM users
ORDER BY username;

-- name: CreateHaiku :execresult
INSERT INTO haikus (from_user_id, to_user_id, line1, line2, line3)
VALUES (?, ?, ?, ?, ?);

-- name: GetReceivedHaikus :many
SELECT 
    h.id,
    h.from_user_id,
    h.to_user_id,
    h.line1,
    h.line2,
    h.line3,
    h.created_at,
    u.username as from_username
FROM haikus h
JOIN users u ON h.from_user_id = u.id
WHERE h.to_user_id = ?
ORDER BY h.created_at DESC;

-- name: GetSentHaikus :many
SELECT 
    h.id,
    h.from_user_id,
    h.to_user_id,
    h.line1,
    h.line2,
    h.line3,
    h.created_at,
    u.username as to_username
FROM haikus h
JOIN users u ON h.to_user_id = u.id
WHERE h.from_user_id = ?
ORDER BY h.created_at DESC;
