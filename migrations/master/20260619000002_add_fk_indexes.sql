-- +goose Up
-- ============================================================
-- ADD MISSING FK INDEXES
-- These FK columns were not covered by previous index migration.
-- ============================================================

CREATE INDEX idx_artikel_id_penulis
  ON artikel (id_penulis);

CREATE INDEX idx_artikel_id_verifikator
  ON artikel (id_verifikator);

CREATE INDEX idx_rujukan_id_faskes
  ON rujukan (id_faskes);

CREATE INDEX idx_user_account_id_pendidikan
  ON user_account (id_pendidikan);

CREATE INDEX idx_user_account_id_pekerjaan
  ON user_account (id_pekerjaan);

CREATE INDEX idx_user_account_id_pendapatan
  ON user_account (id_pendapatan);

-- +goose Down
DROP INDEX IF EXISTS idx_artikel_id_penulis;
DROP INDEX IF EXISTS idx_artikel_id_verifikator;
DROP INDEX IF EXISTS idx_rujukan_id_faskes;
DROP INDEX IF EXISTS idx_user_account_id_pendidikan;
DROP INDEX IF EXISTS idx_user_account_id_pekerjaan;
DROP INDEX IF EXISTS idx_user_account_id_pendapatan;
