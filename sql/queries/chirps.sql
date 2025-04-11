-- name: CreateChirp :one
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
SELECT * from chirps
ORDER BY created_at ASC;

-- name: GetAllChirpsByUserId :many
SELECT * from chirps
where user_id = $1
ORDER BY created_at ASC;

-- name: GetChirpById :one
SELECT * from chirps
where id = $1;

-- name: DeleteChirpById :exec
DELETE from chirps
where id = $1;
