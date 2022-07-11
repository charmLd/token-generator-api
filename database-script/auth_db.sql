create database auth;

use auth;

CREATE TABLE `users` (
  `user_id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
  `hashed_pass` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `salt` varchar(100) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `hasher` varchar(10) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT 'bcrypt' COMMENT 'store password hash type',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `last_login` datetime NOT NULL,
  `role` varchar(45) NOT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=latin1;

LOCK TABLES `users` WRITE;

INSERT INTO `users` VALUES (1,'chamika@gmail.com','$2a$10$TV2uUdffKGjjdiWOFwjXae5857B1qUkN2uKRi2Sjp5khjUiI0sdgW','AbC','bcrypt','2022-07-10 20:19:58','2022-07-10 20:19:58','2022-07-10 18:50:58','admin');

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

LOCK TABLES `auth_tokens` WRITE;

INSERT INTO `auth_tokens` VALUES (_binary 'Fó\ã*S^KW§ÿ¥€\ÇK',1,0,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ0b2tlbl9pZCI6IjQ2ZjNlMzJhLTUzNWUtNGI1Ny1hN2ZmLTlkYTVhNGM3NGIwMiIsInVzZXJfcm9sZSI6ImFkbWluIiwiaWF0IjoxNjU3NDcxODU1LCJleHAiOjE2NTc1NTgyNTV9.Lf1WjXlN01KxOBK9LNXVkSvoOeUIEI-fLY0Du9oq0_o','2022-07-10 22:20:55','2022-07-10 22:20:55'),(_binary '¿Q3\ÆHLÑYj;\Ë',1,0,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ0b2tlbl9pZCI6ImJmODk1MTBiLTMzYzYtNDg0Yy05NmQxLTllNTk2YTNiY2IxYyIsInVzZXJfcm9sZSI6IiIsImlhdCI6MTY1NzQ2NTU3NCwiZXhwIjoxNjU4MDcwMzc0fQ.7jrduyvmtXj8a72BGFoIz3BoEyWerilUl7sgIr_guUU','2022-07-10 20:36:14','2022-07-10 20:36:14'),(_binary '\Ê%\êo¥(BÛ\Ê*\"û1',1,0,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ0b2tlbl9pZCI6ImNhMjVlYTZmLWE1MjgtNDJkYi04OTg2LWNhMmEyMjE4ZmIzMSIsInVzZXJfcm9sZSI6ImFkbWluIiwiaWF0IjoxNjU3NDY1NzM3LCJleHAiOjE2NTgwNzA1Mzd9.M1RmNN4u24H2WMGTE8k8v-_QG8pzc19Cjb2vAS0WnB0','2022-07-10 20:38:57','2022-07-10 20:38:57');

UNLOCK TABLES;


DROP TABLE IF EXISTS `generated_tokens`;

CREATE TABLE `generated_tokens` (
  `gen_token_id` binary(16) NOT NULL,
  `user_id` bigint unsigned NOT NULL,
  `is_blacklisted` tinyint(1) NOT NULL DEFAULT '0',
  `token` text NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `expiry` datetime NOT NULL,
  PRIMARY KEY (`gen_token_id`),
  KEY `gen_tokens_FK_1` (`user_id`),
  CONSTRAINT `gen_tokens_FK_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

LOCK TABLES `generated_tokens` WRITE;
/*!40000 ALTER TABLE `generated_tokens` DISABLE KEYS */;
INSERT INTO `generated_tokens` VALUES (_binary '|ÆC!«;ª\Ãy·\Ì',1,1,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbiI6IjAxODI5MDdjLWM2OGUtNDMyMS04ZWFiLTNiYWFjMzc5YjdjYyIsImlhdCI6MTY1NzQ3MDM4NSwiZXhwIjoxNjU4MDc1MTg1LCJpc19ibGFja2xpc3RlZCI6ZmFsc2V9.jk8e9tHB_vEggZ_6ROxrBPSK_nrApqlt6IGNUaBTtyo','0001-01-01 05:19:24','0001-01-01 05:19:24'),(_binary '\âº\ÙCµix/erŒý',1,1,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbiI6IjhlZTIwNTg2LWJhZDktNDM5Ny1iNTY5LTc4MmY2NTcyYmNmZCIsImlhdCI6MTY1NzQ3MTU2OSwiZXhwIjoxNjU4MDc2MzY5LCJpc19ibGFja2xpc3RlZCI6ZmFsc2V9.PP74JgBwKyIpXQGF_Od4KAyX9o9HwXhkqZVhVTlxKqk','2022-07-10 22:16:09','2022-07-10 22:16:09'),(_binary 'œŠû\'AY6!Ç€\×8(',1,0,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbiI6IjFYVmwiLCJpYXQiOjE2NTc0Njg3MDUsImV4cCI6MTY1ODA3MzUwNSwiaXNfYmxhY2tsaXN0ZWQiOmZhbHNlfQ.csbInX-NpCqf_prndDHy2gmWIGDwPR1MsiyJ_7A1vRI','0001-01-01 05:19:24','0001-01-01 05:19:24'),(_binary 'ò(È#D.±\êØ¿\ì',1,0,'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0b2tlbiI6IjFYVmwiLCJpYXQiOjE2NTc0Njk0NjMsImV4cCI6MTY1ODA3NDI2MywiaXNfYmxhY2tsaXN0ZWQiOmZhbHNlfQ.CtzxVJlQqbvfVSQheaB7JzAf2yOzvpIhcf2uk7PEnDo','0001-01-01 05:19:24','0001-01-01 05:19:24');

UNLOCK TABLES;


