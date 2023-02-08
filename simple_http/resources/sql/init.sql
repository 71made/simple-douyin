CREATE TABLE `user`
(
    `id`             bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
    `username`       varchar(128) NOT NULL DEFAULT '' COMMENT 'Username',
    `password`       varchar(128) NOT NULL DEFAULT '' COMMENT 'Password',
    `follow_count`   bigint NOT NULL DEFAULT 0 COMMENT 'User follow count',
    `follower_count` bigint NOT NULL DEFAULT 0 COMMENT 'User follower count',
    `created_at`     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'User account create time',
    `updated_at`     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'User account update time',
    `deleted_at`     timestamp NULL DEFAULT NULL COMMENT 'User account delete time',
    PRIMARY KEY (`id`),
    KEY              `idx_username` (`username`) COMMENT 'Username index'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='User account table';


CREATE TABLE `user_video`
(
    `id`             bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
    `author_id`      bigint unsigned NOT NULL COMMENT 'Author id',
    `play_uri`       varchar(128) NOT NULL DEFAULT '' COMMENT 'Video resource uri',
    `cover_uri`      varchar(128) NOT NULL DEFAULT '' COMMENT 'Video cover image uri',
    `favorite_count` bigint NOT NULL DEFAULT 0 COMMENT 'Video favorite count',
    `comment_count`  bigint NOT NULL DEFAULT 0 COMMENT 'Video comment count',
    `title`          varchar(128) NOT NULL DEFAULT '' COMMENT 'Video title',
    `created_at`     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'User video create time',
    `updated_at`     timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'User video update time',
    `deleted_at`     timestamp NULL DEFAULT NULL COMMENT 'User video delete time',
    PRIMARY KEY (`id`),
    FOREIGN KEY (`author_id`) REFERENCES `user`(`id`),
    KEY              `idx_author_id_title` (`author_id`, `title`) COMMENT 'Author id and title index',
    KEY              `idx_created_at` (`created_at`) COMMENT 'Video create time index'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='User video table';

CREATE TABLE `favorite`
(
    `id`          bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
    `user_id`     bigint unsigned NOT NULL COMMENT 'User id',
    `video_id`    bigint unsigned NOT NULL COMMENT 'Video id',
    `is_favorite` tinyint unsigned NOT NULL DEFAULT 2 COMMENT 'Favorite: 1-true, 2-false',
    `created_at`  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'User video create time',
    `updated_at`  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'User video update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`user_id`, `video_id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
    FOREIGN KEY (`video_id`) REFERENCES `user_video`(`id`),
    KEY            `idx_author_id_video_id` (`user_id`, `video_id`) COMMENT 'User id and video id index',
    KEY            `idx_updated_at` (`updated_at`) COMMENT 'User favorite update time index'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='User favorite video table';

CREATE TABLE `video_comment`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
    `user_id`    bigint unsigned NOT NULL COMMENT 'User id',
    `video_id`   bigint unsigned NOT NULL COMMENT 'Video id',
    `content`    varchar(255) NOT NULL COMMENT 'Comment content',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Video comment create time',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Video comment update time',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'Video comment delete time',
    PRIMARY KEY (`id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
    FOREIGN KEY (`video_id`) REFERENCES `user_video`(`id`),
    KEY          `idx_video_id` (`video_id`) COMMENT 'Video id index',
    KEY          `idx_created_at` (`created_at`) COMMENT 'Video comment create time index'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Video comment table';

CREATE TABLE `user_relation`
        (
        `id`           bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
        `user_id`      bigint unsigned NOT NULL COMMENT 'User id',
        `follower_id`  bigint unsigned NOT NULL COMMENT 'Follower user id',
        `is_following` tinyint unsigned NOT NULL DEFAULT 2 COMMENT 'Following: 1-true, 2-false',
        `is_friend`    tinyint unsigned NOT NULL DEFAULT 2 COMMENT 'Friend(Already following each other): 1-true, 2-false',
        `created_at`   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'User relation create time',
        `updated_at`   timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'User relation update time',
    PRIMARY KEY (`id`),
    UNIQUE KEY (`user_id`, `follower_id`),
    FOREIGN KEY (`user_id`) REFERENCES `user`(`id`),
    FOREIGN KEY (`follower_id`) REFERENCES `user`(`id`),
    KEY            `idx_user_id_follower_id` (`user_id`, `follower_id`) COMMENT 'User follower ids index',
    KEY            `idx_follower_id_user_id` (`follower_id`, `user_id`) COMMENT 'User follow ids index',
    KEY            `idx_updated_at` (`updated_at`) COMMENT 'Follow relation update time index'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='User relation table';

CREATE TABLE `message`
(
    `id`         bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
    `from_user_id`    bigint unsigned NOT NULL COMMENT 'User id',
    `to_user_id` bigint unsigned NOT NULL COMMENT 'To user id',
    `content`    text NOT NULL COMMENT 'Message content',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Message create time',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Message update time',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'Message delete time',
    PRIMARY KEY (`id`),
    FOREIGN KEY (`from_user_id`) REFERENCES `user`(`id`),
    FOREIGN KEY (`to_user_id`) REFERENCES `user`(`id`),
    KEY            `idx_from_user_id_to_user_id` (`from_user_id`, `to_user_id`) COMMENT 'User send message index',
    KEY            `idx_to_user_id_from_user_id` (`to_user_id`, `from_user_id`) COMMENT 'Message to user index',
    KEY            `idx_updated_at` (`created_at`) COMMENT 'Message update time index'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='User Message table';
