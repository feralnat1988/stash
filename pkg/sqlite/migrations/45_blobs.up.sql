CREATE TABLE `blobs` (
    `checksum` varchar(255) NOT NULL PRIMARY KEY,
    `blob` blob
);

ALTER TABLE `tags` ADD COLUMN `image_checksum` blob REFERENCES `blobs`(`checksum`);
ALTER TABLE `studios` ADD COLUMN `image_checksum` blob REFERENCES `blobs`(`checksum`);
ALTER TABLE `performers` ADD COLUMN `image_checksum` blob REFERENCES `blobs`(`checksum`);

-- ALTER TABLE `scenes` ADD COLUMN `cover_checksum` blob REFERENCES `blobs`(`checksum`);

-- TODO: migrate scenes_cover to cover_checksum - post-migrate