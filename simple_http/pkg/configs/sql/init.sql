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

# CREATE TABLE `note`
# (
#     `id`         bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'PK',
#     `user_id`    int(64) NOT NULL DEFAULT 0 COMMENT 'UserID',
#     `title`      varchar(128) NOT NULL DEFAULT '' COMMENT 'Title',
#     `content`    TEXT NULL COMMENT 'Content',
#     `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'Note create time',
#     `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Note update time',
#     `deleted_at` timestamp NULL DEFAULT NULL COMMENT 'Note delete time',
#     PRIMARY KEY (`id`),
#     KEY          `idx_user_id_title` (`user_id`, `title`) COMMENT 'UserID Title index'
# ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='Note table';