-- +goose Up
ALTER TABLE users
  ADD COLUMN model_player TEXT,
  ADD COLUMN model_deck TEXT,
  ADD COLUMN model_meta TEXT,
  ADD COLUMN le_input_players TEXT,
  ADD COLUMN le_target_players TEXT,
  ADD COLUMN le_input_decks TEXT,
  ADD COLUMN le_target_decks TEXT;

ALTER TABLE users
  DROP COLUMN IF EXISTS decision_tree;

-- +goose Down
ALTER TABLE users
  DROP COLUMN IF EXISTS model_player,
  DROP COLUMN IF EXISTS model_deck,
  DROP COLUMN IF EXISTS model_meta,
  DROP COLUMN IF EXISTS le_input_players,
  DROP COLUMN IF EXISTS le_target_players,
  DROP COLUMN IF EXISTS le_input_decks,
  DROP COLUMN IF EXISTS le_target_decks;

ALTER TABLE users
  ADD COLUMN decision_tree TEXT;
