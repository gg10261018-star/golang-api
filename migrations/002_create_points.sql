CREATE TABLE IF NOT EXISTS `user_points` (
    `email`   VARCHAR(50) NOT NULL,
    `balance` INT         NOT NULL DEFAULT 0,
    `updated` DATETIME    NOT NULL,
    PRIMARY KEY (`email`),
    CONSTRAINT `fk_user_points_email` FOREIGN KEY (`email`) REFERENCES `user` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- request_id is the idempotency key: no matter how many times the same
-- request_id is submitted, the UNIQUE KEY only lets one row through.
-- That's the real safety net; the app-level "check then insert" is just a fast path.
CREATE TABLE IF NOT EXISTS `point_transactions` (
    `id`            BIGINT                    NOT NULL AUTO_INCREMENT,
    `email`         VARCHAR(50)               NOT NULL,
    `request_id`    VARCHAR(64)               NOT NULL,
    `amount`        INT                       NOT NULL,
    `status`        ENUM('SUCCESS','FAILED')  NOT NULL,
    `balance_after` INT                       NOT NULL,
    `message`       VARCHAR(255)              NOT NULL DEFAULT '',
    `created`       DATETIME                  NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_request_id` (`request_id`),
    KEY `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
