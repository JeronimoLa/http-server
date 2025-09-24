-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: AddChirpsToUser :one
INSERT INTO chirps (chirp_id, created_at, updated_at, body, id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

