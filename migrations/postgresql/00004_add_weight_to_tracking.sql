-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE trackings
ADD weight VARCHAR(255) NOT NULL DEFAULT 'Unknown';

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE add_tracking_dialogs
DROP COLUMN weight;