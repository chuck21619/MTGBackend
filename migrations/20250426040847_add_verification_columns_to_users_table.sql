-- +goose Up
ALTER TABLE users 
ADD COLUMN email_verified BOOLEAN DEFAULT FALSE, 
ADD COLUMN verification_token TEXT,
ADD COLUMN verification_token_expires_at TIMESTAMP;

-- +goose Down
ALTER TABLE users 
DROP COLUMN email_verified,
DROP COLUMN verification_token,
DROP COLUMN verification_token_expires_at;
