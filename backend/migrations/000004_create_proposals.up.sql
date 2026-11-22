CREATE TABLE proposals (
  id                BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  uuid              UUID            NOT NULL,
  leader_account_id BIGINT UNSIGNED NOT NULL,
  title             VARCHAR(100)    NOT NULL,
  description       VARCHAR(2000)   NOT NULL DEFAULT '',
  -- Optimistic lock counter; exposed as the ETag and bumped on every write.
  version           INT UNSIGNED    NOT NULL DEFAULT 1,
  created_at        DATETIME(3)     NOT NULL DEFAULT UTC_TIMESTAMP(3),
  -- No ON UPDATE: MariaDB only allows CURRENT_TIMESTAMP there (session time zone),
  -- so the UPDATE statement sets UTC_TIMESTAMP(3) explicitly.
  updated_at        DATETIME(3)     NOT NULL DEFAULT UTC_TIMESTAMP(3),
  PRIMARY KEY (id),
  UNIQUE KEY uq_proposals_uuid (uuid),
  UNIQUE KEY uq_proposals_leader_title (leader_account_id, title),
  CONSTRAINT fk_proposals_leader FOREIGN KEY (leader_account_id)
    REFERENCES accounts (id) ON DELETE RESTRICT,
  CONSTRAINT ck_proposals_title CHECK (CHAR_LENGTH(TRIM(title)) > 0),
  CONSTRAINT ck_proposals_version CHECK (version >= 1)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
