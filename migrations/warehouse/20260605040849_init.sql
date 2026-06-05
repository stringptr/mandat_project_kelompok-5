-- +goose Up

-- ============================================================
-- DATA WAREHOUSE - MONITORING PREVALENSI STUNTING
-- ============================================================

-- ============================================================
-- DIMENSI
-- ============================================================

-- ------------------------------------------------------------
-- DIM_WAKTU
-- Diisi generate untuk rentang tahun yang dibutuhkan
-- ------------------------------------------------------------
CREATE TABLE DIM_WAKTU (
    id_waktu    INT             PRIMARY KEY,  -- format YYYYMMDD, e.g. 20240115
    tanggal     DATE            NOT NULL UNIQUE,
    tahun       INT             NOT NULL CHECK (tahun >= 2000),
    semester    INT             NOT NULL CHECK (semester IN (1, 2)),
    kuartal     INT             NOT NULL CHECK (kuartal BETWEEN 1 AND 4),
    bulan       INT             NOT NULL CHECK (bulan BETWEEN 1 AND 12),
    nama_bulan  VARCHAR(20)     NOT NULL,
    minggu      INT             NOT NULL CHECK (minggu BETWEEN 1 AND 53)
);

-- ------------------------------------------------------------
-- DIM_LOKASI
-- Denormalisasi dari tabel lokasi OLTP (hierarki flatten)
-- ------------------------------------------------------------
CREATE TABLE DIM_LOKASI (
    id_lokasi       SERIAL          PRIMARY KEY,
    id_lokasi_oltp  INT             NOT NULL UNIQUE,   -- FK ke lokasi.id_lokasi tingkat paling bawah
    nama_kelurahan  VARCHAR(255),
    nama_kecamatan  VARCHAR(255),
    nama_kabupaten  VARCHAR(255),
    nama_kota       VARCHAR(255),
    nama_provinsi   VARCHAR(255)
);

-- ------------------------------------------------------------
-- DIM_PASIEN
-- Satu baris per orang (ibu maupun anak)
-- tipe_pasien membedakan konteks baris ini
-- ------------------------------------------------------------
CREATE TYPE tipe_pasien AS ENUM ('Anak', 'Ibu Hamil');
CREATE TYPE kategori_usia AS ENUM ('Balita', 'Anak-Anak', 'Remaja', 'Dewasa');

CREATE TABLE DIM_PASIEN (
    id_dim_pasien       SERIAL          PRIMARY KEY,
    tipe_pasien         tipe_pasien  NOT NULL,
    id_pasien_oltp      INT             NOT NULL,       -- FK ke pasien.id_pasien
    jenis_kelamin       VARCHAR(20),
    usia                INT             CHECK (usia >= 0),
    satuan_usia         VARCHAR(10),                    -- 'Bulan' atau 'Tahun'
    kategori_usia       kategori_usia,                  -- 'Balita', 'Remaja', 'Anak-Anak', 'Dewasa'
    pendidikan          VARCHAR(100),
    pekerjaan           VARCHAR(100),
    kategori_pendapatan VARCHAR(100),
    jumlah_tanggungan   INT             CHECK (jumlah_tanggungan >= 0),
    UNIQUE (id_pasien_oltp, tipe_pasien)
);

-- ------------------------------------------------------------
-- DIM_ANAK
-- Outrigger dari DIM_PASIEN untuk pasien bertipe 'Anak'
-- id_ibu_hamil_oltp nullable — ibu bisa tidak terdaftar
-- ------------------------------------------------------------
CREATE TABLE DIM_ANAK (
    id_dim_pasien           SERIAL       PRIMARY KEY,
    berat_lahir             NUMERIC(5,2) NOT NULL,
    panjang_lahir           NUMERIC(5,2) NOT NULL,
    hubungan_dengan_wali    VARCHAR(50)  NOT NULL,
    id_ibu_hamil_oltp       INT,                        -- NULL jika ibu bukan pasien
    CONSTRAINT fk_dim_anak_pasien
        FOREIGN KEY (id_dim_pasien)
        REFERENCES DIM_PASIEN (id_dim_pasien)
        ON DELETE CASCADE
);

-- ------------------------------------------------------------
-- DIM_IBU_HAMIL
-- Outrigger dari DIM_PASIEN untuk pasien bertipe 'Ibu Hamil'
-- Granular per kehamilan — satu ibu bisa banyak baris
-- PK menggunakan id_dim_pasien + hamil_ke agar satu ibu
-- bisa punya beberapa kehamilan terdaftar
-- ------------------------------------------------------------
CREATE TABLE DIM_IBU_HAMIL (
    id_dim_ibu_hamil    SERIAL          PRIMARY KEY,
    id_dim_pasien       INT             NOT NULL,                    -- FK ke DIM_PASIEN (ibu)
    id_ibu_hamil_oltp   INT             NOT NULL UNIQUE,             -- FK ke ibu_hamil.id_ibu_hamil
    hamil_ke            INT             NOT NULL CHECK (hamil_ke > 0),
    status_kehamilan    VARCHAR(50)     NOT NULL,
    trimester           INT             NOT NULL,                    -- '1/2/3', 'Nifas', dll
    hpht                DATE,
    bulan_mulai_hamil   DATE,
    UNIQUE (id_dim_pasien, hamil_ke),
    CONSTRAINT fk_dim_ibu_hamil_pasien
        FOREIGN KEY (id_dim_pasien)
        REFERENCES DIM_PASIEN (id_dim_pasien)
        ON DELETE CASCADE
);

-- ------------------------------------------------------------
-- DIM_PETUGAS
-- ------------------------------------------------------------
CREATE TYPE peran_petugas AS ENUM ('Kader', 'Bidan', 'Dinas Kesehatan');

CREATE TABLE DIM_PETUGAS (
    id_petugas      SERIAL PRIMARY KEY,
    nama_petugas    VARCHAR(255),
    peran           peran_petugas NOT NULL
);

-- ------------------------------------------------------------
-- DIM_POSYANDU
-- Denormalisasi posyandu
-- ------------------------------------------------------------
CREATE TABLE DIM_POSYANDU (
    id_posyandu         SERIAL          PRIMARY KEY,
    id_posyandu_oltp    INT             NOT NULL UNIQUE,
    nama_posyandu       VARCHAR(255)    NOT NULL,
    id_petugas_bidan    INT             NOT NULL,
    wilayah_kerja       VARCHAR(255)    NOT NULL,
    CONSTRAINT fk_bidan_dim_posyandu
        FOREIGN KEY (id_petugas_bidan)
        REFERENCES DIM_PETUGAS (id_petugas)
        ON DELETE CASCADE
);


-- ============================================================
-- FACT TABLE
-- ============================================================

-- ------------------------------------------------------------
-- FACT_PEMERIKSAAN
-- Grain : satu baris per sesi pemeriksaan (anak maupun ibu)
-- id_dim_kehamilan nullable — hanya terisi untuk ibu hamil
--   yang terdaftar sebagai pasien
-- ------------------------------------------------------------
CREATE TABLE FACT_PEMERIKSAAN (
    id_fact                      BIGSERIAL       PRIMARY KEY,

    -- Degenerate dimension
    id_jadwal_imunisasi_oltp     INT             NOT NULL,   -- referensi ke OLTP

    -- Foreign keys ke dimensi
    id_waktu                     INT             NOT NULL,
    id_lokasi                    INT             NOT NULL,
    id_dim_pasien                INT             NOT NULL,
    id_posyandu                  INT             NOT NULL,
    id_dim_kehamilan             INT,                        -- NULL untuk anak / ibu tidak terdaftar

    -- Konteks pasien
    tipe_pasien                  tipe_pasien  NOT NULL,
    id_petugas                   INT             NOT NULL,

    -- Measures fisik
    berat_badan                  NUMERIC(5,2)    NOT NULL,
    tinggi_badan                 NUMERIC(5,2)    NOT NULL,
    lingkar_kepala               NUMERIC(5,2),
    tekanan_darah                VARCHAR(20),

    -- Measures status klinis
    status_stunting              VARCHAR(50),
    status_gizi                  VARCHAR(50),

    -- Flag biner (1 = ya, 0 = tidak) untuk agregasi cepat
    flag_stunting                SMALLINT        NOT NULL DEFAULT 0 CHECK (flag_stunting IN (0,1)),
    flag_gizi_buruk              SMALLINT        NOT NULL DEFAULT 0 CHECK (flag_gizi_buruk IN (0,1)),
    flag_perlu_rujukan           SMALLINT        NOT NULL DEFAULT 0 CHECK (flag_perlu_rujukan IN (0,1)),
    flag_kek                     SMALLINT        NOT NULL DEFAULT 0 CHECK (flag_kek IN (0,1)),         -- ibu: LILA < 23.5 cm
    flag_anemia                  SMALLINT        NOT NULL DEFAULT 0 CHECK (flag_anemia IN (0,1)),      -- ibu: Hb < 11 g/dL
    flag_risiko_tinggi           SMALLINT        NOT NULL DEFAULT 0 CHECK (flag_risiko_tinggi IN (0,1)),

    -- Audit
    created_at                   TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_fp_waktu
        FOREIGN KEY (id_waktu)      REFERENCES DIM_WAKTU    (id_waktu),
    CONSTRAINT fk_fp_lokasi
        FOREIGN KEY (id_lokasi)     REFERENCES DIM_LOKASI   (id_lokasi),
    CONSTRAINT fk_fp_pasien
        FOREIGN KEY (id_dim_pasien) REFERENCES DIM_PASIEN   (id_dim_pasien),
    CONSTRAINT fk_fp_petugas
        FOREIGN KEY (id_petugas)    REFERENCES DIM_PETUGAS  (id_petugas),
    CONSTRAINT fk_fp_posyandu
        FOREIGN KEY (id_posyandu)   REFERENCES DIM_POSYANDU (id_posyandu),
    CONSTRAINT fk_fp_kehamilan
        FOREIGN KEY (id_dim_kehamilan) REFERENCES DIM_IBU_HAMIL (id_dim_ibu_hamil)
);

-- ------------------------------------------------------------
-- FACT_IMUNISASI
-- Grain : satu baris per jadwal imunisasi
-- ------------------------------------------------------------
CREATE TABLE FACT_IMUNISASI (
    id_fact                 BIGSERIAL       PRIMARY KEY,

    -- Foreign keys ke dimensi
    id_waktu                INT             NOT NULL,
    id_lokasi               INT             NOT NULL,
    id_dim_pasien           INT             NOT NULL,
    id_posyandu             INT             NOT NULL,

    -- Measures deskriptif
    nama_vaksin             VARCHAR(255)    NOT NULL,
    status_imunisasi        VARCHAR(50)     NOT NULL,

    -- Flag biner
    flag_terlaksana         SMALLINT        NOT NULL DEFAULT 0 CHECK (flag_terlaksana IN (0,1)),
    flag_terlambat          SMALLINT        NOT NULL DEFAULT 0 CHECK (flag_terlambat IN (0,1)),

    -- Measures numerik
    hari_keterlambatan      INT             CHECK (hari_keterlambatan >= 0),

    -- Audit
    created_at              TIMESTAMPTZ     NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_fi_waktu
        FOREIGN KEY (id_waktu)      REFERENCES DIM_WAKTU    (id_waktu),
    CONSTRAINT fk_fi_lokasi
        FOREIGN KEY (id_lokasi)     REFERENCES DIM_LOKASI   (id_lokasi),
    CONSTRAINT fk_fi_pasien
        FOREIGN KEY (id_dim_pasien) REFERENCES DIM_PASIEN   (id_dim_pasien),
    CONSTRAINT fk_fi_posyandu
        FOREIGN KEY (id_posyandu)   REFERENCES DIM_POSYANDU (id_posyandu)
);


-- ============================================================
-- INDEX — mempercepat query OLAP umum
-- ============================================================

-- FACT_PEMERIKSAAN
CREATE INDEX idx_fp_waktu        ON FACT_PEMERIKSAAN (id_waktu);
CREATE INDEX idx_fp_lokasi       ON FACT_PEMERIKSAAN (id_lokasi);
CREATE INDEX idx_fp_pasien       ON FACT_PEMERIKSAAN (id_dim_pasien);
CREATE INDEX idx_fp_posyandu     ON FACT_PEMERIKSAAN (id_posyandu);
CREATE INDEX idx_fp_tipe_pasien  ON FACT_PEMERIKSAAN (tipe_pasien);
CREATE INDEX idx_fp_kehamilan    ON FACT_PEMERIKSAAN (id_dim_kehamilan);
CREATE INDEX idx_fp_flags        ON FACT_PEMERIKSAAN (flag_stunting, flag_kek, flag_anemia);

-- FACT_IMUNISASI
CREATE INDEX idx_fi_waktu        ON FACT_IMUNISASI (id_waktu);
CREATE INDEX idx_fi_lokasi       ON FACT_IMUNISASI (id_lokasi);
CREATE INDEX idx_fi_pasien       ON FACT_IMUNISASI (id_dim_pasien);
CREATE INDEX idx_fi_posyandu     ON FACT_IMUNISASI (id_posyandu);
CREATE INDEX idx_fi_vaksin       ON FACT_IMUNISASI (nama_vaksin);

-- DIM_PASIEN
CREATE INDEX idx_dp_tipe         ON DIM_PASIEN (tipe_pasien);
CREATE INDEX idx_dp_oltp         ON DIM_PASIEN (id_pasien_oltp);

-- DIM_LOKASI
CREATE INDEX idx_dl_kabupaten    ON DIM_LOKASI (nama_kabupaten);
CREATE INDEX idx_dl_kecamatan    ON DIM_LOKASI (nama_kecamatan);

-- DIM_IBU_HAMIL
CREATE INDEX idx_dih_pasien      ON DIM_IBU_HAMIL (id_dim_pasien);
CREATE INDEX idx_dih_trimester   ON DIM_IBU_HAMIL (trimester);


-- ============================================================
-- COMMENT
-- ============================================================
COMMENT ON COLUMN FACT_PEMERIKSAAN.flag_kek                 IS 'Kurang Energi Kronis: LILA ibu < 23.5 cm. Derive dari status_gizi saat ETL.';
COMMENT ON COLUMN FACT_PEMERIKSAAN.flag_anemia              IS 'Anemia: Hb ibu < 11 g/dL. Derive dari catatan klinis saat ETL.';
COMMENT ON COLUMN FACT_PEMERIKSAAN.id_dim_kehamilan         IS 'NULL untuk tipe_pasien=Anak, atau ibu yang tidak terdaftar sebagai pasien.';
COMMENT ON COLUMN FACT_PEMERIKSAAN.id_jadwal_imunisasi_oltp IS 'Degenerate dimension — referensi ke jadwal_imunisasi.id_imunisasi di OLTP.';
COMMENT ON COLUMN DIM_ANAK.id_ibu_hamil_oltp                IS 'NULL jika ibu tidak terdaftar sebagai pasien di sistem.';
COMMENT ON COLUMN DIM_PASIEN.satuan_usia                    IS 'Bulan (untuk anak < 5 tahun) atau Tahun.';
COMMENT ON COLUMN DIM_WAKTU.id_waktu                        IS 'Format YYYYMMDD, e.g. 20240115 untuk 15 Januari 2024.';

-- +goose Down

-- ============================================================
-- DROP INDEXES
-- ============================================================

-- FACT_PEMERIKSAAN indexes
DROP INDEX IF EXISTS idx_fp_waktu;
DROP INDEX IF EXISTS idx_fp_lokasi;
DROP INDEX IF EXISTS idx_fp_pasien;
DROP INDEX IF EXISTS idx_fp_posyandu;
DROP INDEX IF EXISTS idx_fp_tipe_pasien;
DROP INDEX IF EXISTS idx_fp_kehamilan;
DROP INDEX IF EXISTS idx_fp_flags;

-- FACT_IMUNISASI indexes
DROP INDEX IF EXISTS idx_fi_waktu;
DROP INDEX IF EXISTS idx_fi_lokasi;
DROP INDEX IF EXISTS idx_fi_pasien;
DROP INDEX IF EXISTS idx_fi_posyandu;
DROP INDEX IF EXISTS idx_fi_vaksin;

-- DIM_PASIEN indexes
DROP INDEX IF EXISTS idx_dp_tipe;
DROP INDEX IF EXISTS idx_dp_oltp;

-- DIM_LOKASI indexes
DROP INDEX IF EXISTS idx_dl_kabupaten;
DROP INDEX IF EXISTS idx_dl_kecamatan;

-- DIM_IBU_HAMIL indexes
DROP INDEX IF EXISTS idx_dih_pasien;
DROP INDEX IF EXISTS idx_dih_trimester;


-- ============================================================
-- DROP TABLES
-- ============================================================

-- Fact Tables
DROP TABLE IF EXISTS FACT_IMUNISASI;
DROP TABLE IF EXISTS FACT_PEMERIKSAAN;

-- Dimension Tables (with outriggers after their parent)
DROP TABLE IF EXISTS DIM_ANAK;
DROP TABLE IF EXISTS DIM_IBU_HAMIL;
DROP TABLE IF EXISTS DIM_POSYANDU;
DROP TABLE IF EXISTS DIM_PETUGAS;
DROP TABLE IF EXISTS DIM_PASIEN;
DROP TABLE IF EXISTS DIM_LOKASI;
DROP TABLE IF EXISTS DIM_WAKTU;


-- ============================================================
-- DROP CUSTOM TYPES
-- ============================================================
DROP TYPE IF EXISTS peran_petugas;
DROP TYPE IF EXISTS kategori_usia;
DROP TYPE IF EXISTS tipe_pasien;
