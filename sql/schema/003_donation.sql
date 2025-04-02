-- +goose Up
CREATE TABLE donation(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    amount INTEGER NOT NULL,
    user_id UUID NOT NULL,
    non_profit_name TEXT NOT NULL,
    
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_non_profit_name FOREIGN KEY (non_profit_name) REFERENCES non_profit(name) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE donation;
