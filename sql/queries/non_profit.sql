-- name: CreateNonProfit :one
INSERT INTO non_profit(id, created_at, updated_at, name, address, email, owner_id)
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

-- name: DeleteNonProfits :exec
DELETE FROM non_profit; 


-- name: GetNonProfByName :one
SELECT * FROM non_profit
WHERE name = $1;
