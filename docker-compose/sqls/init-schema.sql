CREATE DATABASE IF NOT EXISTS `posts_list_demo` character set utf8mb4;

USE `posts_list_demo`;

CREATE TABLE IF NOT EXISTS `posts` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `content` longtext NOT NULL,
  `like_count` int(11) NOT NULL DEFAULT '0',
  `publish_time` datetime NOT NULL,
  `publisher_id` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
