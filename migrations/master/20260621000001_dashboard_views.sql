-- +goose Up
-- +goose StatementBegin

-- View 1: Distribusi status gizi (latest pemeriksaan per pasien)
CREATE VIEW v_dashboard_distribusi_gizi AS
WITH latest AS (
    SELECT DISTINCT ON (ji.id_pasien)
        hp.status_gizi,
        ji.id_pasien
    FROM hasil_pemeriksaan hp
    JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
    JOIN pasien p ON ji.id_pasien = p.id_pasien
    WHERE p.is_deleted = FALSE
    ORDER BY ji.id_pasien, hp.created_at DESC
)
SELECT status_gizi, COUNT(*)::int AS jumlah
FROM latest
GROUP BY status_gizi;

-- View 2: Tren stunting bulanan (distinct pasien with stunting per month)
CREATE VIEW v_dashboard_tren_stunting AS
SELECT
    TO_CHAR(hp.created_at, 'YYYY-MM') AS bulan,
    COUNT(DISTINCT ji.id_pasien)::int AS jumlah
FROM hasil_pemeriksaan hp
JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
WHERE hp.status_stunting IN ('Stunting', 'Stunting Berat')
GROUP BY TO_CHAR(hp.created_at, 'YYYY-MM')
ORDER BY bulan;

-- View 3: Stunting per wilayah (kabupaten/kota level)
CREATE VIEW v_dashboard_stunting_per_wilayah AS
WITH latest AS (
    SELECT DISTINCT ON (ji.id_pasien)
        hp.status_stunting,
        ji.id_pasien
    FROM hasil_pemeriksaan hp
    JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
    JOIN pasien p ON ji.id_pasien = p.id_pasien
    WHERE p.is_deleted = FALSE
    ORDER BY ji.id_pasien, hp.created_at DESC
),
pasien_lokasi AS (
    SELECT
        p.id_pasien,
        COALESCE(l_kab.nama_lokasi, 'Tidak Diketahui') AS nama_wilayah
    FROM pasien p
    JOIN posyandu pos ON p.id_posyandu = pos.id_posyandu
    JOIN lokasi l_kel ON pos.id_lokasi = l_kel.id_lokasi
    LEFT JOIN lokasi l_kec ON l_kel.bagian_dari = l_kec.id_lokasi AND l_kec.tipe_lokasi = 'Kecamatan'
    LEFT JOIN lokasi l_kab ON l_kec.bagian_dari = l_kab.id_lokasi AND l_kab.tipe_lokasi = 'Kabupaten'
    WHERE p.is_deleted = FALSE
)
SELECT
    pl.nama_wilayah,
    COUNT(DISTINCT pl.id_pasien)::int AS total_balita,
    COUNT(DISTINCT CASE WHEN l.status_stunting IN ('Stunting', 'Stunting Berat') THEN pl.id_pasien END)::int AS jumlah_kasus,
    ROUND(
        COUNT(DISTINCT CASE WHEN l.status_stunting IN ('Stunting', 'Stunting Berat') THEN pl.id_pasien END)::numeric * 100.0
        / NULLIF(COUNT(DISTINCT pl.id_pasien), 0), 1
    ) AS prevalensi,
    CASE
        WHEN COUNT(DISTINCT pl.id_pasien) > 0
            AND COUNT(DISTINCT CASE WHEN l.status_stunting IN ('Stunting', 'Stunting Berat') THEN pl.id_pasien END)::numeric * 100.0
                / COUNT(DISTINCT pl.id_pasien) >= 20
        THEN 'tinggi'
        WHEN COUNT(DISTINCT pl.id_pasien) > 0
            AND COUNT(DISTINCT CASE WHEN l.status_stunting IN ('Stunting', 'Stunting Berat') THEN pl.id_pasien END)::numeric * 100.0
                / COUNT(DISTINCT pl.id_pasien) >= 10
        THEN 'sedang'
        ELSE 'rendah'
    END AS level
FROM pasien_lokasi pl
LEFT JOIN latest l ON pl.id_pasien = l.id_pasien
GROUP BY pl.nama_wilayah;

-- View 4: Kehadiran posyandu bulanan (distinct pasien diperiksa per bulan)
CREATE VIEW v_dashboard_kehadiran_bulanan AS
SELECT
    TO_CHAR(hp.created_at, 'YYYY-MM') AS bulan,
    COUNT(DISTINCT ji.id_pasien)::int AS jumlah
FROM hasil_pemeriksaan hp
JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
GROUP BY TO_CHAR(hp.created_at, 'YYYY-MM')
ORDER BY bulan;

-- View 0: Dashboard stats (total pasien, perlu verifikasi, tindak lanjut, kasus stunting, cakupan)
CREATE VIEW v_dashboard_stats AS
WITH latest AS (
    SELECT DISTINCT ON (ji.id_pasien)
        hp.status_stunting,
        ji.id_pasien
    FROM hasil_pemeriksaan hp
    JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
    JOIN pasien p ON ji.id_pasien = p.id_pasien
    WHERE p.is_deleted = FALSE
    ORDER BY ji.id_pasien, hp.created_at DESC
)
SELECT
    (SELECT COUNT(*)::int FROM pasien WHERE is_deleted = FALSE) AS total_pasien,
    (SELECT COUNT(*)::int
     FROM hasil_pemeriksaan hp
     WHERE hp.created_at >= CURRENT_DATE - INTERVAL '7 days') AS perlu_verifikasi,
    (SELECT COUNT(*)::int
     FROM tindak_lanjut tl
     WHERE tl.status_pasien != 'Selesai Pemantauan') AS tindak_lanjut,
    (SELECT COUNT(*)::int
     FROM latest
     WHERE status_stunting IN ('Stunting', 'Stunting Berat')) AS kasus_stunting,
    (SELECT COUNT(*)::int
     FROM jadwal_imunisasi
     WHERE status_imunisasi = 'Belum'
       AND tanggal_jadwal >= CURRENT_DATE) AS jadwal_posyandu,
    (SELECT COUNT(*)::int
     FROM pasien p JOIN anak a ON p.id_pasien = a.id_pasien
     WHERE p.is_deleted = FALSE) AS total_balita,
    (SELECT ROUND(
        COUNT(*) FILTER (WHERE status_imunisasi = 'Sudah')::numeric * 100.0
        / NULLIF(COUNT(*), 0), 1
     )
     FROM jadwal_imunisasi) AS cakupan_persentase;

-- View 4.5: Jadwal imunisasi terdekat
CREATE VIEW v_dashboard_jadwal_terdekat AS
SELECT
    ji.id_imunisasi AS id,
    ji.nama_vaksin,
    ji.tanggal_jadwal,
    COALESCE(a.nama_anak, ua.nama) AS nama_pasien
FROM jadwal_imunisasi ji
JOIN pasien p ON ji.id_pasien = p.id_pasien
JOIN user_account ua ON p.id_pasien = ua.id_user
LEFT JOIN anak a ON p.id_pasien = a.id_pasien
WHERE ji.status_imunisasi = 'Belum'
  AND ji.tanggal_jadwal >= CURRENT_DATE
  AND p.is_deleted = FALSE
ORDER BY ji.tanggal_jadwal ASC
LIMIT 10;

-- View 5: Public stats (tanpa auth, kasus stunting = latest per pasien)
CREATE VIEW v_public_stats AS
WITH latest AS (
    SELECT DISTINCT ON (ji.id_pasien)
        hp.status_stunting,
        ji.id_pasien
    FROM hasil_pemeriksaan hp
    JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
    JOIN pasien p ON ji.id_pasien = p.id_pasien
    WHERE p.is_deleted = FALSE
    ORDER BY ji.id_pasien, hp.created_at DESC
)
SELECT
    (SELECT COUNT(*)::int FROM pasien WHERE is_deleted = FALSE) AS total_pasien,
    (SELECT COUNT(*)::int FROM pasien p JOIN anak a ON p.id_pasien = a.id_pasien WHERE p.is_deleted = FALSE) AS balita_dipantau,
    (SELECT COUNT(*)::int FROM latest WHERE status_stunting IN ('Stunting', 'Stunting Berat')) AS kasus_stunting,
    (SELECT COUNT(*)::int FROM artikel WHERE status_artikel = 'Dipublikasikan') AS total_artikel;

-- View 6: Riwayat pemeriksaan per pasien
CREATE VIEW v_riwayat_pemeriksaan AS
SELECT
    ji.id_pasien,
    hp.created_at AS tanggal,
    hp.berat_badan::numeric AS berat_badan,
    hp.tinggi_badan::numeric AS tinggi_badan,
    hp.status_gizi,
    hp.catatan,
    ua.nama AS petugas
FROM hasil_pemeriksaan hp
JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
JOIN user_account ua ON hp.id_petugas_input = ua.id_user
ORDER BY ji.id_pasien, hp.created_at DESC;

-- View 7: Tumbuh kembang (latest per pasien per bulan)
CREATE VIEW v_tumbuh_kembang AS
SELECT DISTINCT ON (ji.id_pasien, TO_CHAR(hp.created_at, 'YYYY-MM'))
    ji.id_pasien,
    TO_CHAR(hp.created_at, 'YYYY-MM') AS bulan,
    hp.berat_badan::numeric AS berat_badan,
    hp.tinggi_badan::numeric AS tinggi_badan
FROM hasil_pemeriksaan hp
JOIN jadwal_imunisasi ji ON hp.id_jadwal_imunisasi = ji.id_imunisasi
ORDER BY ji.id_pasien, TO_CHAR(hp.created_at, 'YYYY-MM'), hp.created_at DESC;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP VIEW IF EXISTS v_tumbuh_kembang;
DROP VIEW IF EXISTS v_riwayat_pemeriksaan;
DROP VIEW IF EXISTS v_public_stats;
DROP VIEW IF EXISTS v_dashboard_kehadiran_bulanan;
DROP VIEW IF EXISTS v_dashboard_stunting_per_wilayah;
DROP VIEW IF EXISTS v_dashboard_tren_stunting;
DROP VIEW IF EXISTS v_dashboard_jadwal_terdekat;
DROP VIEW IF EXISTS v_dashboard_distribusi_gizi;
DROP VIEW IF EXISTS v_dashboard_stats;

-- +goose StatementEnd
