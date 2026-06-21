-- +goose Up
-- ============================================================
-- ADD INDEXES for query performance
-- Master DB currently has ZERO explicit indexes beyond PK/UNIQUE.
-- With millions of rows, FK joins and filtered queries need indexes.
-- ============================================================

-- ============================================================
-- HIGH PRIORITY — hot path queries (dashboard, lists, filters)
-- ============================================================

CREATE INDEX idx_user_account_lokasi
  ON user_account (id_lokasi);

CREATE INDEX idx_user_account_deleted
  ON user_account (is_deleted);

CREATE INDEX idx_user_account_status_verifikasi
  ON user_account (status_verifikasi);

CREATE INDEX idx_user_account_lokasi_deleted
  ON user_account (id_lokasi, is_deleted);

CREATE INDEX idx_pasien_posyandu
  ON pasien (id_posyandu);

CREATE INDEX idx_jadwal_imunisasi_pasien
  ON jadwal_imunisasi (id_pasien);

CREATE INDEX idx_jadwal_imunisasi_pasien_tanggal
  ON jadwal_imunisasi (id_pasien, tanggal_jadwal);

CREATE INDEX idx_hasil_pemeriksaan_jadwal
  ON hasil_pemeriksaan (id_jadwal_imunisasi);

CREATE INDEX idx_hasil_pemeriksaan_petugas
  ON hasil_pemeriksaan (id_petugas_input);

CREATE INDEX idx_hasil_pemeriksaan_petugas_tanggal
  ON hasil_pemeriksaan (id_petugas_input, created_at);

CREATE INDEX idx_notifikasi_user
  ON notifikasi (id_user);

CREATE INDEX idx_notifikasi_user_baca
  ON notifikasi (id_user, status_baca);

CREATE INDEX idx_ibu_hamil_pasien
  ON ibu_hamil (id_pasien);

-- ============================================================
-- MEDIUM PRIORITY — common FK joins
-- ============================================================

CREATE INDEX idx_lokasi_bagian_dari
  ON lokasi (bagian_dari);

CREATE INDEX idx_bidan_wilayah_kerja
  ON bidan (wilayah_kerja);

CREATE INDEX idx_posyandu_lokasi
  ON posyandu (id_lokasi);

CREATE INDEX idx_posyandu_bidan
  ON posyandu (id_bidan);

CREATE INDEX idx_kader_posyandu_posyandu
  ON kader_posyandu (id_posyandu);

CREATE INDEX idx_fasilitas_kesehatan_lokasi
  ON fasilitas_kesehatan (id_lokasi);

CREATE INDEX idx_anak_ibu_hamil
  ON anak (id_ibu_hamil);

CREATE INDEX idx_anak_wali
  ON anak (id_wali);

CREATE INDEX idx_tindak_lanjut_bidan
  ON tindak_lanjut (id_bidan);

CREATE INDEX idx_user_session_user
  ON user_session (id_user);

CREATE INDEX idx_audit_log_user
  ON audit_log (id_user);

CREATE INDEX idx_audit_log_session
  ON audit_log (id_user_session);

-- ============================================================
-- FILTER / SORT indexes for paginated list endpoints
-- ============================================================

CREATE INDEX idx_artikel_status
  ON artikel (status_artikel);

CREATE INDEX idx_artikel_tanggal_publish
  ON artikel (tanggal_publish DESC);

CREATE INDEX idx_rujukan_status
  ON rujukan (status_rujukan);

CREATE INDEX idx_tindak_lanjut_status
  ON tindak_lanjut (status_pasien);

CREATE INDEX idx_hasil_pemeriksaan_status_gizi
  ON hasil_pemeriksaan (status_gizi);

-- +goose Down
-- ============================================================
-- ROLLBACK — drop all indexes
-- ============================================================

DROP INDEX IF EXISTS idx_user_account_lokasi;
DROP INDEX IF EXISTS idx_user_account_deleted;
DROP INDEX IF EXISTS idx_user_account_status_verifikasi;
DROP INDEX IF EXISTS idx_user_account_lokasi_deleted;
DROP INDEX IF EXISTS idx_pasien_posyandu;
DROP INDEX IF EXISTS idx_jadwal_imunisasi_pasien;
DROP INDEX IF EXISTS idx_jadwal_imunisasi_pasien_tanggal;
DROP INDEX IF EXISTS idx_hasil_pemeriksaan_jadwal;
DROP INDEX IF EXISTS idx_hasil_pemeriksaan_petugas;
DROP INDEX IF EXISTS idx_hasil_pemeriksaan_petugas_tanggal;
DROP INDEX IF EXISTS idx_notifikasi_user;
DROP INDEX IF EXISTS idx_notifikasi_user_baca;
DROP INDEX IF EXISTS idx_ibu_hamil_pasien;
DROP INDEX IF EXISTS idx_lokasi_bagian_dari;
DROP INDEX IF EXISTS idx_bidan_wilayah_kerja;
DROP INDEX IF EXISTS idx_posyandu_lokasi;
DROP INDEX IF EXISTS idx_posyandu_bidan;
DROP INDEX IF EXISTS idx_kader_posyandu_posyandu;
DROP INDEX IF EXISTS idx_fasilitas_kesehatan_lokasi;
DROP INDEX IF EXISTS idx_anak_ibu_hamil;
DROP INDEX IF EXISTS idx_anak_wali;
DROP INDEX IF EXISTS idx_tindak_lanjut_bidan;
DROP INDEX IF EXISTS idx_user_session_user;
DROP INDEX IF EXISTS idx_audit_log_user;
DROP INDEX IF EXISTS idx_audit_log_session;
DROP INDEX IF EXISTS idx_artikel_status;
DROP INDEX IF EXISTS idx_artikel_tanggal_publish;
DROP INDEX IF EXISTS idx_rujukan_status;
DROP INDEX IF EXISTS idx_tindak_lanjut_status;
DROP INDEX IF EXISTS idx_hasil_pemeriksaan_status_gizi;
