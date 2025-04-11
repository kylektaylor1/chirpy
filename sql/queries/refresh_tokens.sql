-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id, expires_at) 
VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: GetRefreshTokenByToken :one
SELECT * from refresh_tokens
where token = $1;

-- name: SetRefreshTokenRevoked :exec
UPDATE refresh_tokens
set revoked_at = NOW(), updated_at = NOW()
where token = $1;
