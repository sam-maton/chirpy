-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET updated_at = NOW(), email = $2, hashed_password = $3
WHERE id = $1
RETURNING id, email, created_at, updated_at, is_chirpy_red;

-- name: UpgradeUser :one
UPDATE users
SET updated_at = NOW(), is_chirpy_red = true
WHERE id = $1
RETURNING id, email, updated_at, is_chirpy_red;