-- +goose Up
ALTER TABLE users
    ALTER COLUMN decision_tree TYPE BYTEA
    USING decode(encode(decision_tree::bytea, 'hex'), 'hex');

-- +goose Down
ALTER TABLE users
    ALTER COLUMN decision_tree TYPE TEXT;
