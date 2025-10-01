-- name: AddChirpsToUser :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps;

-- name: GetSingleChirp :one
SELECT * FROM chirps
WHERE id = $1; 

-- name: DeleteSingleChirp :exec
DELETE FROM chirps
WHERE id = $1;

-- name: GetAuthorChirps :many
SELECT * FROM chirps
WHERE user_id = $1; 

