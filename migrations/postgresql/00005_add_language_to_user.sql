-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE users
ADD language VARCHAR(255) NOT NULL DEFAULT 'ru';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE users
DROP COLUMN language;