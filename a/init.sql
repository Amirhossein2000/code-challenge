create database if not exists `appDB`;

USE `appDB`;

CREATE TABLE IF NOT EXISTS `orders` (
    `id` int(11) NOT NULL,
    `price` int NOT NULL,
    `title` varchar(30) NOT NULL,
    PRIMARY KEY  (`id`)
);
