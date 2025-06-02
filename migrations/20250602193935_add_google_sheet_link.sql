-- +goose Up
ALTER TABLE users ADD COLUMN google_sheet TEXT;

-- +goose Down
ALTER TABLE users DROP COLUMN google_sheet;
