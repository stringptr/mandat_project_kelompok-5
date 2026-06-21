-- +goose Up
-- +goose StatementBegin

CREATE MATERIALIZED VIEW mv_ibu_hamil_stats AS
SELECT
    COUNT(DISTINCT ih.id_pasien)::int AS total_ibu_hamil,
    COUNT(DISTINCT CASE WHEN ih.status_kehamilan = 'Trimester 1' THEN ih.id_pasien END)::int AS trimester_1,
    COUNT(DISTINCT CASE WHEN ih.status_kehamilan = 'Trimester 2' THEN ih.id_pasien END)::int AS trimester_2,
    COUNT(DISTINCT CASE WHEN ih.status_kehamilan = 'Trimester 3' THEN ih.id_pasien END)::int AS trimester_3,
    COUNT(DISTINCT CASE WHEN ih.status_kehamilan = 'Melahirkan' THEN ih.id_pasien END)::int AS melahirkan,
    COUNT(DISTINCT CASE WHEN ih.status_kehamilan = 'Nifas' THEN ih.id_pasien END)::int AS nifas,
    COUNT(DISTINCT CASE WHEN ih.status_kehamilan = 'Keguguran' THEN ih.id_pasien END)::int AS keguguran
FROM ibu_hamil ih
JOIN pasien p ON ih.id_pasien = p.id_pasien
WHERE p.is_deleted = FALSE
  AND ih.is_deleted = FALSE;

CREATE MATERIALIZED VIEW mv_ibu_hamil_per_wilayah AS
WITH pasien_lokasi AS (
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
    COUNT(DISTINCT pl.id_pasien)::int AS total_ibu_hamil,
    COUNT(DISTINCT CASE WHEN ih.status_kehamilan = 'Trimester 1' THEN pl.id_pasien END)::int AS trimester_1,
    COUNT(DISTINCT CASE WHEN ih.status_kehamilan = 'Trimester 2' THEN pl.id_pasien END)::int AS trimester_2,
    COUNT(DISTINCT CASE WHEN ih.status_kehamilan = 'Trimester 3' THEN pl.id_pasien END)::int AS trimester_3,
    COUNT(DISTINCT CASE WHEN ih.status_kehamilan = 'Melahirkan' THEN pl.id_pasien END)::int AS melahirkan,
    COUNT(DISTINCT CASE WHEN ih.status_kehamilan = 'Nifas' THEN pl.id_pasien END)::int AS nifas,
    COUNT(DISTINCT CASE WHEN ih.status_kehamilan = 'Keguguran' THEN pl.id_pasien END)::int AS keguguran
FROM pasien_lokasi pl
JOIN ibu_hamil ih ON pl.id_pasien = ih.id_pasien
WHERE ih.is_deleted = FALSE
GROUP BY pl.nama_wilayah
ORDER BY total_ibu_hamil DESC;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP MATERIALIZED VIEW IF EXISTS mv_ibu_hamil_per_wilayah;
DROP MATERIALIZED VIEW IF EXISTS mv_ibu_hamil_stats;

-- +goose StatementEnd
