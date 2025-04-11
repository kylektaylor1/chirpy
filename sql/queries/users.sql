-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * from users
where email = $1;

-- name: UpdateUser :exec
UPDATE users
set updated_at = NOW(), email = $1, hashed_password = $2
where id = $3;

-- name: GetUserById :one
SELECT * from users
where id = $1;

-- name: UpgradeUserById :exec
UPDATE users
set updated_at = NOW(), is_chirpy_red = true 
where id = $1;


