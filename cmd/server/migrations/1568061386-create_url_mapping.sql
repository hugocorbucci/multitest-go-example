-- Migration: create_url_mapping
-- Created at: 2019-09-09 17:36:26
-- ====  UP  ====

START TRANSACTION;
  CREATE TABLE IF NOT EXISTS url_mapping (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT NOT NULL,
    url TEXT COLLATE utf8mb4_unicode_ci NOT NULL,
    short_url VARCHAR(12) NOT NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
  ) DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

  CREATE UNIQUE INDEX idx_url_mapping_short_url
    ON url_mapping (short_url);
COMMIT;

-- ==== DOWN ====

START TRANSACTION;
  DROP INDEX idx_url_mapping_short_url ON url_mapping;
  DROP TABLE IF EXISTS url_mapping;
COMMIT;
