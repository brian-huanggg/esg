ALTER TABLE sessions
  ADD COLUMN csrf_token_hash BINARY(32) NOT NULL AFTER token_hash;