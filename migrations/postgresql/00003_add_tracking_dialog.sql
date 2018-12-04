
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE add_tracking_dialogs (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    created TIMESTAMP NOT NULL,
    modified TIMESTAMP NOT NULL,
    step INTEGER NOT NULL DEFAULT 0,
    user_id INTEGER NOT NULL UNIQUE,
    future_tracking_name VARCHAR(255) NOT NULL DEFAULT '',
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE add_tracking_dialogs
DROP FOREIGN KEY user_id;
DROP TABLE add_tracking_dialogs;