CREATE DATABASE IF NOT EXISTS `gochat`;
USE `gochat`;

CREATE TABLE IF NOT EXISTS `users` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(45) NOT NULL,
  `password_hash` varchar(60) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `username_UNIQUE` (`username`)
) ;

CREATE TABLE IF NOT EXISTS `chats` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(45) NOT NULL,
  `is_private` tinyint NOT NULL DEFAULT '0',
  `creator` int unsigned NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `fk_chats_creator_idx` (`creator`),
  CONSTRAINT `fk_chats_creator` FOREIGN KEY (`creator`) REFERENCES `users` (`id`)
) ;

CREATE TABLE IF NOT EXISTS `chat_participants` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int unsigned NOT NULL,
  `chat_id` int unsigned NOT NULL,
  `accepted_invite` tinyint NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `fk_chat_participants_user_id_idx` (`user_id`),
  KEY `fk_chat_participants_chat_id_idx` (`chat_id`),
  CONSTRAINT `fk_chat_participants_chat_id` FOREIGN KEY (`chat_id`) REFERENCES `chats` (`id`),
  CONSTRAINT `fk_chat_participants_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ;


CREATE TABLE IF NOT EXISTS `message` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `chat_id` int unsigned NOT NULL,
  `user_id` int unsigned NOT NULL,
  `message` varchar(200) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  KEY `fk_message_chat_id_idx` (`chat_id`),
  KEY `fk_message_user_id_idx` (`user_id`),
  CONSTRAINT `fk_message_chat_id` FOREIGN KEY (`chat_id`) REFERENCES `chats` (`id`),
  CONSTRAINT `fk_message_user_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ;

