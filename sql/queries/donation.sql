-- name: CreateDonation :one
INSERT INTO donation(id,  created_at, updated_at ,amount, user_id, non_profit_name)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    $3
)
    RETURNING *;


-- name: ShowDonations :many
SELECT non_profit_name as name, amount
FROM donation
JOIN non_profit
ON non_profit.name = donation.non_profit_name
WHERE user_id = $1;

