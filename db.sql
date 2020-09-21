DROP DATABASE IF EXISTS thescore;

CREATE DATABASE thescore;

USE thescore;

DROP TABLE IF EXISTS rushingstats;

CREATE TABLE `rushingstats`(
    `stat_id` int(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    `player` varchar(128) NOT NULL,
    `yards` int(10) NOT NULL,
    `longest` int(10) NOT NULL,
    `touchdowns` int(10) NOT NULL,
    `misc` JSON,
    INDEX `idx_player` (`player`),
    INDEX `idx_yards` (`yards`),
    INDEX `idx_longest` (`longest`),
    INDEX `idx_touchdowns` (`touchdowns`)
) ENGINE=InnoDB CHARACTER SET utf8;

FLUSH TABLES;

