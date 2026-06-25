-- 1. TRUNCATE semua tabel
TRUNCATE TABLE
  lokasi,
  pendidikan,
  kategori_pendapatan,
  pekerjaan,
  user_account,
  dinas_kesehatan,
  bidan,
  fasilitas_kesehatan,
  posyandu,
  kader_posyandu,
  pasien,
  ibu_hamil,
  anak,
  jadwal_imunisasi,
  artikel,
  hasil_pemeriksaan,
  tindak_lanjut,
  rujukan,
  notifikasi,
  audit_log,
  user_session
RESTART IDENTITY CASCADE;

-- 2. Bulk insert (tabel tanpa is_deleted)
\copy lokasi FROM '/data-csv/lokasi.csv' CSV HEADER;
\copy pendidikan FROM '/data-csv/pendidikan.csv' CSV HEADER;
\copy kategori_pendapatan FROM '/data-csv/kategori_pendapatan.csv' CSV HEADER;
\copy pekerjaan FROM '/data-csv/pekerjaan.csv' CSV HEADER;

-- tabel dengan is_deleted — explicit column list
\copy user_account (id_user,email,password,no_hp,status_verifikasi,nama,nik,jenis_kelamin,tanggal_lahir,id_lokasi,id_pendidikan,id_pekerjaan,id_pendapatan,jumlah_tanggungan,akun_ke,created_at,updated_at) FROM '/data-csv/user_account.csv' CSV HEADER;

\copy dinas_kesehatan (id_user,created_at,updated_at) FROM '/data-csv/dinas_kesehatan.csv' CSV HEADER;
\copy bidan (id_user,no_str,wilayah_kerja,created_at,updated_at) FROM '/data-csv/bidan.csv' CSV HEADER;
\copy fasilitas_kesehatan (id_faskes,nama_faskes,tipe_faskes,id_lokasi,created_at,updated_at) FROM '/data-csv/fasilitas_kesehatan.csv' CSV HEADER;
\copy posyandu (id_posyandu,nama_posyandu,id_lokasi,id_bidan,created_at,updated_at) FROM '/data-csv/posyandu.csv' CSV HEADER;
\copy kader_posyandu (id_user,no_sk,id_posyandu,created_at,updated_at) FROM '/data-csv/kader_posyandu.csv' CSV HEADER;

\copy pasien (id_pasien,id_posyandu,created_at,updated_at) FROM '/data-csv/pasien.csv' CSV HEADER;
\copy ibu_hamil (id_ibu_hamil,id_pasien,hamil_ke,bulan_mulai_hamil,hpht,status_kehamilan,created_at,updated_at) FROM '/data-csv/ibu_hamil.csv' CSV HEADER;

-- anak: bypass FK karena seed data mungkin inkonsisten
SET session_replication_role = replica;
\copy anak (id_pasien,id_ibu_hamil,id_wali,nama_anak,berat_lahir,panjang_lahir,hubungan_dengan_wali,created_at,updated_at) FROM '/data-csv/anak.csv' CSV HEADER;
SET session_replication_role = default;

\copy jadwal_imunisasi FROM '/data-csv/jadwal_imunisasi.csv' CSV HEADER;
\copy artikel FROM '/data-csv/artikel.csv' CSV HEADER;

\copy hasil_pemeriksaan FROM '/data-csv/hasil_pemeriksaan.csv' CSV HEADER;
\copy tindak_lanjut FROM '/data-csv/tindak_lanjut.csv' CSV HEADER;
\copy rujukan FROM '/data-csv/rujukan.csv' CSV HEADER;

\copy notifikasi (id_notifikasi,id_user,judul,pesan,tipe_notifikasi,status_baca,tanggal_kirim) FROM '/data-csv/notifikasi.csv' CSV HEADER;

\copy user_session FROM '/data-csv/user_session.csv' CSV HEADER;

-- audit_log: staging table untuk fix single-quote JSON
CREATE TEMP TABLE audit_log_staging (
    id_log INT, tipe_aktor TEXT, id_user TEXT, id_user_session TEXT, tipe_aktivitas TEXT, berhasil TEXT,
    endpoint TEXT, table_name TEXT, record_id TEXT, old_value TEXT, new_value TEXT,
    detail TEXT, ip_address TEXT, user_agent TEXT, waktu_aktivitas TEXT
);

\copy audit_log_staging FROM '/data-csv/audit_log.csv' CSV HEADER;

INSERT INTO audit_log (id_log, tipe_aktor, id_user, id_user_session, tipe_aktivitas, berhasil, endpoint, table_name, record_id, old_value, new_value, detail, ip_address, user_agent, waktu_aktivitas)
SELECT
    id_log,
    tipe_aktor::tipe_aktor,
    NULLIF(id_user, '')::INT,
    NULLIF(id_user_session, '')::UUID,
    tipe_aktivitas::tipe_aktivitas,
    berhasil::BOOLEAN,
    NULLIF(endpoint, ''),
    NULLIF(table_name, ''),
    NULLIF(record_id, ''),
    REPLACE(REPLACE(NULLIF(old_value, ''), '''', '"'), 'None', 'null')::JSONB,
    REPLACE(REPLACE(NULLIF(new_value, ''), '''', '"'), 'None', 'null')::JSONB,
    NULLIF(detail, ''),
    NULLIF(ip_address, '')::INET,
    NULLIF(user_agent, ''),
    waktu_aktivitas::TIMESTAMPTZ
FROM audit_log_staging;

DROP TABLE audit_log_staging;
