CREATE DATABASE IF NOT EXISTS `web_index`;

CREATE TABLE IF NOT EXISTS `web_keys` (
    `id` bigint(11) NOT NULL AUTO_INCREMENT,
    `key` varchar(255) NOT NULL,
    `count` mediumint(8)  NOT NULL,
    `url` varchar(255) NOT NULL,
    `last_crawled` TIMESTAMP NOT NULL,
    PRIMARY KEY(`id`)
);
