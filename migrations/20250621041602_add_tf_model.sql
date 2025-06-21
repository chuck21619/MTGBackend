-- +goose Up
ALTER TABLE users
ADD COLUMN tf_model BYTEA,
ADD COLUMN tf_players BYTEA,
ADD COLUMN tf_decks BYTEA;

-- +goose Down
ALTER TABLE users
DROP COLUMN tf_model,
DROP COLUMN tf_players,
DROP COLUMN tf_decks;
