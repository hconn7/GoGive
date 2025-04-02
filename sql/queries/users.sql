-- name: CreateUser :one
INSERT INTO users(id, created_at, updated_at, email, hashed_password, is_owner, non_profit_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3,
    $4
    
)
    RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users; 

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: UpdateIsOwner :exec
UPDATE users
SET 
    is_owner= true,
    updated_at = NOW()
WHERE id = $1;
