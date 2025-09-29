
-- name: AddRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, email, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    $2,
    $3,
    $4,
    $5,
    NULL    
)
RETURNING *;

-- name: GetRefreshTokenByEmail :one
SELECT * FROM refresh_tokens 
where email = $1 AND revoked_at is NULL;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET updated_at = $1, revoked_at = $2
WHERE token = $3;


-- name: GetUserByRefreshToken :one
SELECT user_id, revoked_at FROM refresh_tokens 
where token = $1;