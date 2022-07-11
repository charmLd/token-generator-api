CREATE DATABASE  IF NOT EXISTS `auth` ;
USE `auth`;
DROP TABLE IF EXISTS `users`;

CREATE TABLE `users` (
                         `user_id` bigint unsigned NOT NULL AUTO_INCREMENT,
                         `email` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
                         `hashed_pass` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
                         `salt` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
                         `hasher` varchar(10) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT 'bcrypt' COMMENT 'store password hash type',
                         `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                         `last_login` datetime DEFAULT NULL,
                         `role` varchar(45) NOT NULL,
                         PRIMARY KEY (`user_id`),
                         UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=latin1;

LOCK TABLES `users` WRITE;
INSERT INTO `users` VALUES (1,'chamika@gmail.com','$2a$10$TV2uUdffKGjjdiWOFwjXae5857B1qUkN2uKRi2Sjp5khjUiI0sdgW','AbC','bcrypt','2022-07-10 20:19:58','2022-07-10 20:19:58','2022-07-11 17:32:41','admin'),(2,'testclient@gmail.com','$2a$10$4NM8RNSMgieAnyKfTS9clOKT4ThVLebfVuKQvT4TaicsONAudzvLy','cAb','bcrypt','2022-07-11 21:28:09','2022-07-11 21:28:09','2022-07-11 17:39:30','user');

UNLOCK TABLES;

DROP TABLE IF EXISTS `auth_tokens`;

CREATE TABLE `auth_tokens` (
                               `token_id` binary(16) NOT NULL,
                               `user_id` bigint unsigned NOT NULL,
                               `is_blacklisted` tinyint NOT NULL,
                               `auth_token` text NOT NULL,
                               `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                               `expiry` datetime NOT NULL,
                               PRIMARY KEY (`token_id`),
                               KEY `login_tokens_FK_1` (`user_id`),
                               CONSTRAINT `login_tokens_FK_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

DROP TABLE IF EXISTS `generated_tokens`;

CREATE TABLE `generated_tokens` (
                                    `id` int NOT NULL AUTO_INCREMENT,
                                    `gen_token_id` varchar(12) NOT NULL,
                                    `user_id` bigint unsigned NOT NULL,
                                    `is_blacklisted` tinyint(1) NOT NULL DEFAULT '0',
                                    `token` varchar(250) NOT NULL,
                                    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                    `expiry` datetime NOT NULL,
                                    PRIMARY KEY (`id`),
                                    UNIQUE KEY `gen_token_user_id_jwt_UNIQUE` (`gen_token_id`,`user_id`,`token`),
                                    KEY `gen_tokens_FK_1` (`user_id`),
                                    CONSTRAINT `gen_tokens_FK_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=latin1;