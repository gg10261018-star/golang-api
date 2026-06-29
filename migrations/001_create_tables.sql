CREATE TABLE IF NOT EXISTS `user` (
    `email`    VARCHAR(50)  NOT NULL,
    `password` VARCHAR(255) NOT NULL,
    `created`  DATETIME     NOT NULL,
    `updated`  DATETIME     NOT NULL,
    PRIMARY KEY (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `weather` (
    `id`      BIGINT          NOT NULL AUTO_INCREMENT,
    `city`    VARCHAR(50)     NOT NULL,
    `min_t`   DECIMAL(4,1)    NOT NULL,
    `max_t`   DECIMAL(4,1)    NOT NULL,
    `period`  ENUM('AM','PM') NOT NULL,
    `created` DATETIME        NOT NULL,
    `updated` DATETIME        NOT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uq_city_period` (`city`, `period`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
