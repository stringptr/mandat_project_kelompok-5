-- +goose Up
-- +goose StatementBegin

-- Critical: stunting filter on pemeriksaan (DISTINCT ON + WHERE status_stunting)
CREATE INDEX idx_hasil_pemeriksaan_stunting_created
  ON hasil_pemeriksaan (status_stunting, created_at DESC);

-- Compound: help DISTINCT ON pattern through jadwal_imunisasi join
CREATE INDEX idx_hasil_pemeriksaan_jadwal_created
  ON hasil_pemeriksaan (id_jadwal_imunisasi, created_at DESC);

-- Help count queries scanning by date range
CREATE INDEX idx_hasil_pemeriksaan_created
  ON hasil_pemeriksaan (created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_hasil_pemeriksaan_stunting_created;
DROP INDEX IF EXISTS idx_hasil_pemeriksaan_jadwal_created;
DROP INDEX IF EXISTS idx_hasil_pemeriksaan_created;

-- +goose StatementEnd
