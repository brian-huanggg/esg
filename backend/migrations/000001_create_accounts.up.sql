CREATE TABLE accounts (
  id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  uuid          UUID            NOT NULL,
  employee_id   CHAR(6)         NOT NULL,
  name          VARCHAR(100)    NOT NULL,
  password_hash VARCHAR(255)    NOT NULL,
  created_at    DATETIME(3)     NOT NULL DEFAULT UTC_TIMESTAMP(3),
  PRIMARY KEY (id),
  UNIQUE KEY uq_accounts_uuid (uuid),
  UNIQUE KEY uq_accounts_employee_id (employee_id),
  CONSTRAINT ck_accounts_employee_id CHECK (employee_id REGEXP '^[0-9]{6}$'),
  CONSTRAINT ck_accounts_name CHECK (CHAR_LENGTH(TRIM(name)) > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
