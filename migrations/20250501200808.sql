-- Create "chats" table
CREATE TABLE `chats` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `telegram_id` integer NOT NULL, `title` text NULL);
-- Create index "chats_telegram_id_key" to table: "chats"
CREATE UNIQUE INDEX `chats_telegram_id_key` ON `chats` (`telegram_id`);
-- Create "gathers" table
CREATE TABLE `gathers` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `active` bool NOT NULL DEFAULT true, `created_at` datetime NOT NULL, `closed_at` datetime NULL, `chat_gathers` integer NULL, CONSTRAINT `gathers_chats_gathers` FOREIGN KEY (`chat_gathers`) REFERENCES `chats` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create "memberships" table
CREATE TABLE `memberships` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `is_admin` bool NOT NULL DEFAULT false, `chat_memberships` integer NULL, `user_memberships` integer NULL, CONSTRAINT `memberships_users_memberships` FOREIGN KEY (`user_memberships`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT `memberships_chats_memberships` FOREIGN KEY (`chat_memberships`) REFERENCES `chats` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create "users" table
CREATE TABLE `users` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `telegram_id` integer NOT NULL, `username` text NULL, `first_name` text NULL, `last_name` text NULL, `is_owner` bool NOT NULL DEFAULT false);
-- Create index "users_telegram_id_key" to table: "users"
CREATE UNIQUE INDEX `users_telegram_id_key` ON `users` (`telegram_id`);
-- Create "votes" table
CREATE TABLE `votes` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `count` integer NOT NULL, `voted_at` datetime NOT NULL, `gather_votes` integer NULL, `user_votes` integer NULL, CONSTRAINT `votes_users_votes` FOREIGN KEY (`user_votes`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL, CONSTRAINT `votes_gathers_votes` FOREIGN KEY (`gather_votes`) REFERENCES `gathers` (`id`) ON UPDATE NO ACTION ON DELETE SET NULL);
