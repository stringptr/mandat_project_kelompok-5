package dashboard

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"

	dashboardDomain "github.com/stringptr/SiGizi/backend/internal/domain/dashboard"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetDashboardStats(ctx context.Context) (*dashboardDomain.DashboardStatsRow, error) {
	row := &dashboardDomain.DashboardStatsRow{}
	err := r.db.QueryRow(ctx, `SELECT total_pasien, perlu_verifikasi, tindak_lanjut, kasus_stunting, jadwal_posyandu, total_balita, cakupan_persentase FROM mv_dashboard_stats`).
		Scan(&row.TotalPasien, &row.PerluVerifikasi, &row.TindakLanjut, &row.KasusStunting, &row.JadwalPosyandu, &row.TotalBalita, &row.CakupanPersentase)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (r *Repo) GetDistribusiGizi(ctx context.Context) ([]dashboardDomain.DistribusiGiziRow, error) {
	rows, err := r.db.Query(ctx, `SELECT status_gizi, jumlah FROM mv_dashboard_distribusi_gizi ORDER BY jumlah DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dashboardDomain.DistribusiGiziRow
	for rows.Next() {
		var row dashboardDomain.DistribusiGiziRow
		if err := rows.Scan(&row.StatusGizi, &row.Jumlah); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

func (r *Repo) GetTrenStunting(ctx context.Context) ([]dashboardDomain.TrenRow, error) {
	rows, err := r.db.Query(ctx, `SELECT bulan, jumlah FROM mv_dashboard_tren_stunting ORDER BY bulan`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dashboardDomain.TrenRow
	for rows.Next() {
		var row dashboardDomain.TrenRow
		if err := rows.Scan(&row.Bulan, &row.Jumlah); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

func (r *Repo) GetStuntingPerWilayah(ctx context.Context) ([]dashboardDomain.StuntingWilayahRow, error) {
	rows, err := r.db.Query(ctx, `SELECT nama_wilayah, prevalensi, jumlah_kasus, total_balita, level FROM mv_dashboard_stunting_per_wilayah ORDER BY jumlah_kasus DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dashboardDomain.StuntingWilayahRow
	for rows.Next() {
		var row dashboardDomain.StuntingWilayahRow
		if err := rows.Scan(&row.NamaWilayah, &row.Prevalensi, &row.JumlahKasus, &row.TotalBalita, &row.Level); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

func (r *Repo) GetKehadiranBulanan(ctx context.Context) ([]dashboardDomain.TrenRow, error) {
	rows, err := r.db.Query(ctx, `SELECT bulan, jumlah FROM mv_dashboard_kehadiran_bulanan ORDER BY bulan`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dashboardDomain.TrenRow
	for rows.Next() {
		var row dashboardDomain.TrenRow
		if err := rows.Scan(&row.Bulan, &row.Jumlah); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

func (r *Repo) GetPublicStats(ctx context.Context) (*dashboardDomain.PublicStatsRow, error) {
	row := &dashboardDomain.PublicStatsRow{}
	err := r.db.QueryRow(ctx, `SELECT total_pasien, balita_dipantau, kasus_stunting, total_artikel FROM mv_public_stats`).
		Scan(&row.TotalPasien, &row.BalitaDipantau, &row.KasusStunting, &row.TotalArtikel)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (r *Repo) GetRiwayat(ctx context.Context, idPasien int32) ([]dashboardDomain.RiwayatRow, error) {
	rows, err := r.db.Query(ctx, `SELECT tanggal::text, berat_badan, tinggi_badan, status_gizi, catatan, petugas FROM mv_riwayat_pemeriksaan WHERE id_pasien = $1 ORDER BY tanggal DESC`, idPasien)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dashboardDomain.RiwayatRow
	for rows.Next() {
		var row dashboardDomain.RiwayatRow
		var berat, tinggi float64
		if err := rows.Scan(&row.Tanggal, &berat, &tinggi, &row.StatusGizi, &row.Catatan, &row.Petugas); err != nil {
			return nil, err
		}
		row.BeratBadan = berat
		row.TinggiBadan = tinggi
		results = append(results, row)
	}
	return results, rows.Err()
}

func (r *Repo) GetTumbuhKembang(ctx context.Context, idPasien int32) ([]dashboardDomain.TumbuhKembangRow, error) {
	rows, err := r.db.Query(ctx, `SELECT bulan, berat_badan, tinggi_badan FROM mv_tumbuh_kembang WHERE id_pasien = $1 ORDER BY bulan`, idPasien)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dashboardDomain.TumbuhKembangRow
	for rows.Next() {
		var row dashboardDomain.TumbuhKembangRow
		var berat, tinggi float64
		if err := rows.Scan(&row.Bulan, &berat, &tinggi); err != nil {
			return nil, err
		}
		row.BeratBadan = berat
		row.TinggiBadan = tinggi
		results = append(results, row)
	}
	return results, rows.Err()
}

func (r *Repo) GetJadwalTerdekat(ctx context.Context) ([]dashboardDomain.JadwalTerdekatRow, error) {
	rows, err := r.db.Query(ctx, `SELECT id, nama_vaksin, tanggal_jadwal::text, nama_pasien FROM mv_dashboard_jadwal_terdekat`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dashboardDomain.JadwalTerdekatRow
	for rows.Next() {
		var row dashboardDomain.JadwalTerdekatRow
		if err := rows.Scan(&row.ID, &row.NamaVaksin, &row.TanggalJadwal, &row.NamaPasien); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

func (r *Repo) GetIbuHamilStats(ctx context.Context) (*dashboardDomain.IbuHamilStatsRow, error) {
	row := &dashboardDomain.IbuHamilStatsRow{}
	err := r.db.QueryRow(ctx, `SELECT total_ibu_hamil, trimester_1, trimester_2, trimester_3, melahirkan, nifas, keguguran FROM mv_ibu_hamil_stats`).
		Scan(&row.TotalIbuHamil, &row.Trimester1, &row.Trimester2, &row.Trimester3, &row.Melahirkan, &row.Nifas, &row.Keguguran)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func (r *Repo) GetIbuHamilPerWilayah(ctx context.Context) ([]dashboardDomain.IbuHamilWilayahRow, error) {
	rows, err := r.db.Query(ctx, `SELECT nama_wilayah, total_ibu_hamil, trimester_1, trimester_2, trimester_3, melahirkan, nifas, keguguran FROM mv_ibu_hamil_per_wilayah ORDER BY total_ibu_hamil DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []dashboardDomain.IbuHamilWilayahRow
	for rows.Next() {
		var row dashboardDomain.IbuHamilWilayahRow
		if err := rows.Scan(&row.NamaWilayah, &row.TotalIbuHamil, &row.Trimester1, &row.Trimester2, &row.Trimester3, &row.Melahirkan, &row.Nifas, &row.Keguguran); err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// RefreshMaterializedViews refreshes all dashboard materialized views.
// Each view is refreshed with CONCURRENTLY to avoid locking out readers.
// Errors are intentionally ignored so a partial failure doesn't break callers.
func (r *Repo) RefreshMaterializedViews(ctx context.Context) error {
	views := []string{
		"mv_dashboard_stats",
		"mv_dashboard_distribusi_gizi",
		"mv_dashboard_tren_stunting",
		"mv_dashboard_kehadiran_bulanan",
		"mv_dashboard_stunting_per_wilayah",
		"mv_public_stats",
		"mv_riwayat_pemeriksaan",
		"mv_tumbuh_kembang",
		"mv_dashboard_jadwal_terdekat",
		"mv_ibu_hamil_stats",
		"mv_ibu_hamil_per_wilayah",
	}
	for _, v := range views {
		// CONCURRENTLY requires a unique index on the MV; if it fails, fall back.
		_, err := r.db.Exec(ctx, "REFRESH MATERIALIZED VIEW CONCURRENTLY "+v)
		if err != nil {
			// Fall back to regular refresh (blocks reads for a moment, but safe).
			r.db.Exec(ctx, "REFRESH MATERIALIZED VIEW "+v) //nolint:errcheck
		}
	}
	return nil
}

var _ dashboardDomain.Repo = (*Repo)(nil)

func (r *Repo) GetPosyanduByKaderID(ctx context.Context, idKader int32) (int32, error) {
	var id int32
	err := r.db.QueryRow(ctx, `SELECT id_posyandu FROM kader_posyandu WHERE id_user = $1 AND is_deleted = false`, idKader).Scan(&id)
	return id, err
}

func (r *Repo) GetSemuaPemeriksaan(ctx context.Context, page, perPage int, idBidan, idPosyandu int32) ([]dashboardDomain.PemeriksaanRow, int, error) {
	where := ""
	args := []interface{}{}
	argIdx := 1

	if idPosyandu > 0 {
		where = " WHERE p.id_posyandu = $" + strconv.Itoa(argIdx)
		args = append(args, idPosyandu)
		argIdx++
	}
	if idBidan > 0 {
		if where == "" {
			where = " WHERE pos.id_bidan = $" + strconv.Itoa(argIdx)
		} else {
			where += " AND pos.id_bidan = $" + strconv.Itoa(argIdx)
		}
		args = append(args, idBidan)
		argIdx++
	}

	countFrom := ` FROM hasil_pemeriksaan hp
		JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
		JOIN pasien p ON p.id_pasien = ji.id_pasien AND p.is_deleted = false `
	if idBidan > 0 || idPosyandu > 0 {
		countFrom += ` JOIN posyandu pos ON pos.id_posyandu = p.id_posyandu `
	}

	var count struct{ Count int }
	err := r.db.QueryRow(ctx, `SELECT COUNT(*)`+countFrom+where, args...).Scan(&count.Count)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	limitArg := argIdx
	offsetArg := argIdx + 1
	allArgs := append(args, perPage, offset)

	sql := `
		WITH page_ids AS (
			SELECT hp.id_hasil_pemeriksaan
			` + countFrom + where + `
			ORDER BY hp.created_at DESC
			LIMIT $` + strconv.Itoa(limitArg) + ` OFFSET $` + strconv.Itoa(offsetArg) + `
		)
		SELECT hp.id_hasil_pemeriksaan, ji.id_imunisasi, COALESCE(ji.nama_vaksin, '') AS nama_vaksin,
		       COALESCE(a.nama_anak, ua.nama) AS nama_pasien,
		       hp.berat_badan::float8, hp.tinggi_badan::float8,
		       hp.lingkar_kepala::float8, hp.tekanan_darah,
		       hp.status_stunting::text, hp.status_gizi::text,
		       hp.catatan, hp.created_at::text AS tanggal, petugas.nama AS petugas
		FROM hasil_pemeriksaan hp
		JOIN page_ids f ON f.id_hasil_pemeriksaan = hp.id_hasil_pemeriksaan
		JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
		JOIN pasien p ON p.id_pasien = ji.id_pasien
		JOIN user_account ua ON ua.id_user = p.id_pasien
		JOIN user_account petugas ON petugas.id_user = hp.id_petugas_input
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien
		ORDER BY hp.created_at DESC
	`

	rows, err := r.db.Query(ctx, sql, allArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []dashboardDomain.PemeriksaanRow
	for rows.Next() {
		var row dashboardDomain.PemeriksaanRow
		if err := rows.Scan(&row.IDHasilPemeriksaan, &row.IDJadwalImunisasi, &row.NamaVaksin,
			&row.NamaPasien, &row.BeratBadan, &row.TinggiBadan,
			&row.LingkarKepala, &row.TekananDarah, &row.StatusStunting, &row.StatusGizi,
			&row.Catatan, &row.Tanggal, &row.Petugas); err != nil {
			return nil, 0, err
		}
		results = append(results, row)
	}
	return results, count.Count, rows.Err()
}
