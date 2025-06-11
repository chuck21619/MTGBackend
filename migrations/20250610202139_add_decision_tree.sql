-- +goose Up
ALTER TABLE users ADD COLUMN decision_tree TEXT;

-- +goose Down
ALTER TABLE users DROP COLUMN decision_tree;
