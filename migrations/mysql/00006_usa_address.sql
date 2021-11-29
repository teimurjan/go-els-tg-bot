
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE usaAddresses (
    `id` BIGSERIAL NOT NULL PRIMARY KEY,
    `street` VARCHAR(255) NOT NULL DEFAULT '',
    `city` VARCHAR(255) NOT NULL DEFAULT '',
    `state` VARCHAR(255) NOT NULL DEFAULT '',
    `zip` VARCHAR(255) NOT NULL DEFAULT '',
    `phone` VARCHAR(255) NOT NULL DEFAULT '',
    `created` TIMESTAMP NOT NULL,
    `modified` TIMESTAMP NOT NULL
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE usaAddresses;
