-- +goose Up

-- ============================================================
-- ADD UNIQUE CONSTRAINTS FOR INCREMENTAL UPSERT
-- ============================================================

-- DIM_PETUGAS: upsert by nama_petugas + peran
ALTER TABLE DIM_PETUGAS ADD CONSTRAINT uq_petugas_nama_peran UNIQUE (nama_petugas, peran);

-- FACT_IMUNISASI: add OLTP reference column for upsert
ALTER TABLE FACT_IMUNISASI ADD COLUMN id_imunisasi_oltp INT;
ALTER TABLE FACT_IMUNISASI ADD CONSTRAINT uq_fi_imunisasi_oltp UNIQUE (id_imunisasi_oltp);

-- FACT_PEMERIKSAAN: add OLTP reference column for upsert
ALTER TABLE FACT_PEMERIKSAAN ADD COLUMN id_hasil_pemeriksaan_oltp INT;
ALTER TABLE FACT_PEMERIKSAAN ADD CONSTRAINT uq_fp_hasil_pemeriksaan_oltp UNIQUE (id_hasil_pemeriksaan_oltp);


-- +goose Down

ALTER TABLE FACT_PEMERIKSAAN DROP CONSTRAINT IF EXISTS uq_fp_hasil_pemeriksaan_oltp;
ALTER TABLE FACT_PEMERIKSAAN DROP COLUMN IF EXISTS id_hasil_pemeriksaan_oltp;

ALTER TABLE FACT_IMUNISASI DROP CONSTRAINT IF EXISTS uq_fi_imunisasi_oltp;
ALTER TABLE FACT_IMUNISASI DROP COLUMN IF EXISTS id_imunisasi_oltp;

ALTER TABLE DIM_PETUGAS DROP CONSTRAINT IF EXISTS uq_petugas_nama_peran;
