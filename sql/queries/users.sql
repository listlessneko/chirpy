-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET 
  email = COALESCE(sqlc.narg('email'), email),
  hashed_password = COALESCE(sqlc.narg('hashed_password'), hashed_password)
WHERE id = $1
RETURNING *;
