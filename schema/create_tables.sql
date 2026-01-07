-- 判断表是否存在，不存在则创建
CREATE TABLE IF NOT EXISTS `discount_rule`
(
    `id`            BIGINT(20) NOT NULL AUTO_INCREMENT,
    `name`          VARCHAR(128) NOT NULL DEFAULT '',
    `description`   VARCHAR(512) NOT NULL DEFAULT '',
    `discount_rate` INT(10) NOT NULL DEFAULT 0,
    `status`        INT(10) NOT NULL DEFAULT 1,
    `created_at`    BIGINT(20) NOT NULL DEFAULT 0,
    `updated_at`    BIGINT(20) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `user_discount` (
                                               `id` BIGINT(20) NOT NULL AUTO_INCREMENT,
    `user_id` BIGINT(20) NOT NULL DEFAULT 0,
    `model_id` INT(10) NOT NULL DEFAULT 0,
    `rule_id` BIGINT(20) NOT NULL DEFAULT 0,
    `discount_rate` INT(10) NOT NULL DEFAULT 0,
    `created_at` BIGINT(20) NOT NULL DEFAULT 0,
    `updated_at` BIGINT(20) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
