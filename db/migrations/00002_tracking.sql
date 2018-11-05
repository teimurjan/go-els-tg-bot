
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `trackings` (
    `id` INTEGER NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `created` DATETIME NOT NULL,
    `modified` DATETIME NOT NULL,
    `name` VARCHAR(255) NOT NULL,
    `value` VARCHAR(255) NOT NULL,
    `status` VARCHAR(255) NOT NULL,
    `user_id` INTEGER NOT NULL,
    FOREIGN KEY (`user_id`) REFERENCES users(`id`)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `trackings`
DROP FOREIGN KEY `user_id`;
DROP TABLE `trackings`;