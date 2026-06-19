package tindaklanjut

import (
	"context"
	"fmt"
	"time"

	tindaklanjutDomain "github.com/stringptr/SiGizi/backend/internal/domain/tindaklanjut"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/enum"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
	. "github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/table"

	"github.com/go-jet/jet/v2/pgxV5"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (r *Repo) GetPasienByID(ctx context.Context, idPasien int32) (*model.Pasien, error) {
	var results []*model.Pasien
	stmt := SELECT(Pasien.AllColumns).
		FROM(Pasien).
		WHERE(Pasien.IDPasien.EQ(Int32(idPasien)).AND(Pasien.IsDeleted.EQ(Bool(false))))
	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}

var statusRujukanExpr = map[model.StatusRujukan]StringExpression{
	model.StatusRujukan_Diajukan: enum.StatusRujukan.Diajukan,
	model.StatusRujukan_Diproses: enum.StatusRujukan.Diproses,
	model.StatusRujukan_Diterima: enum.StatusRujukan.Diterima,
	model.StatusRujukan_Ditolak:  enum.StatusRujukan.Ditolak,
	model.StatusRujukan_Selesai:  enum.StatusRujukan.Selesai,
}

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetPasienTindakLanjut(ctx context.Context, page int, perPage int) ([]*tindaklanjutDomain.PasienTindakLanjutJoinRow, int, error) {
	offset := int64((page - 1) * perPage)

	latestHpSubquery := `
		LEFT JOIN LATERAL (
			SELECT
				hp2.id_hasil_pemeriksaan,
				hp2.status_gizi,
				hp2.status_stunting,
				hp2.catatan,
				hp2.created_at
			FROM hasil_pemeriksaan hp2
			JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp2.id_jadwal_imunisasi
			WHERE ji.id_pasien = p.id_pasien
			ORDER BY hp2.created_at DESC
			LIMIT 1
		) hp ON true
	`

	fromWhere := `
		FROM pasien p
		LEFT JOIN user_account ua ON ua.id_user = p.id_pasien AND ua.is_deleted = false
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
	` + latestHpSubquery + `
		LEFT JOIN tindak_lanjut tl ON tl.id_hasil_pemeriksaan = hp.id_hasil_pemeriksaan
		WHERE p.is_deleted = false
			AND hp.status_gizi IS NOT NULL
			AND hp.status_gizi != 'Gizi Baik'
			AND tl.id_tindak_lanjut IS NULL
	`

	var countResult struct{ Count int64 }
	err := pgxV5.Query(ctx, RawStatement("SELECT COUNT(*)"+fromWhere), r.db, &countResult)
	if err != nil {
		return nil, 0, err
	}

	dataSQL := `
		SELECT
			p.id_pasien,
			COALESCE(ua.nama, a.nama_anak, '') AS nama_pasien,
			hp.status_gizi::text AS status_gizi,
			CASE
				WHEN hp.status_gizi IN ('Gizi Buruk', 'Gizi Kurang')
					OR hp.status_stunting IN ('Stunting', 'Stunting Berat', 'Berisiko Stunting')
				THEN 'Perlu Rujukan'
				ELSE 'Dalam Pemantauan'
			END AS status_pasien,
			hp.created_at::text AS tanggal_pemeriksaan
	` + fromWhere + `
		ORDER BY hp.created_at DESC
		OFFSET $1 LIMIT $2
	`

	pgxRows, err := r.db.Query(ctx, dataSQL, offset, perPage)
	if err != nil {
		return nil, 0, err
	}
	defer pgxRows.Close()

	var rows []*tindaklanjutDomain.PasienTindakLanjutJoinRow
	for pgxRows.Next() {
		var row tindaklanjutDomain.PasienTindakLanjutJoinRow
		err := pgxRows.Scan(&row.IDPasien, &row.NamaPasien, &row.StatusGizi, &row.StatusPasien, &row.TanggalPemeriksaan)
		if err != nil {
			return nil, 0, err
		}
		rows = append(rows, &row)
	}

	return rows, int(countResult.Count), nil
}

func (r *Repo) GetDetailPasienByID(ctx context.Context, idPasien int32) (*tindaklanjutDomain.DetailPasienJoinRow, error) {
	sql := `
		SELECT
			p.id_pasien,
			COALESCE(ua.nama, '') AS nama_pasien,
			CASE WHEN ua.tanggal_lahir IS NOT NULL
				THEN EXTRACT(YEAR FROM AGE(ua.tanggal_lahir))::int || ' Tahun'
				ELSE ''
			END AS usia,
			COALESCE(hp.status_gizi::text, '') AS status_gizi,
			COALESCE(hp.status_stunting::text, '') AS status_stunting,
			COALESCE(hp.catatan, '') AS catatan
		FROM pasien p
		LEFT JOIN user_account ua ON ua.id_user = p.id_pasien AND ua.is_deleted = false
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
		LEFT JOIN LATERAL (
			SELECT hp2.id_hasil_pemeriksaan, hp2.status_gizi, hp2.status_stunting, hp2.catatan
			FROM hasil_pemeriksaan hp2
			JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp2.id_jadwal_imunisasi
			WHERE ji.id_pasien = p.id_pasien
			ORDER BY hp2.created_at DESC
			LIMIT 1
		) hp ON true
		WHERE p.id_pasien = $1 AND p.is_deleted = false
		LIMIT 1
	`

	pgxRows, err := r.db.Query(ctx, sql, idPasien)
	if err != nil {
		return nil, err
	}
	defer pgxRows.Close()

	if !pgxRows.Next() {
		return nil, nil
	}

	var row tindaklanjutDomain.DetailPasienJoinRow
	err = pgxRows.Scan(&row.IDPasien, &row.NamaPasien, &row.Usia, &row.StatusGizi, &row.StatusStunting, &row.Catatan)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) GetRiwayatPemeriksaan(ctx context.Context, idPasien int32) ([]*tindaklanjutDomain.RiwayatPemeriksaanJoinRow, error) {
	sql := `
		SELECT
			hp.created_at::text AS tanggal,
			hp.berat_badan::float8 AS berat_badan,
			hp.tinggi_badan::float8 AS tinggi_badan
		FROM hasil_pemeriksaan hp
		JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
		WHERE ji.id_pasien = $1
		ORDER BY hp.created_at DESC
	`

	pgxRows, err := r.db.Query(ctx, sql, idPasien)
	if err != nil {
		return nil, err
	}
	defer pgxRows.Close()

	var rows []*tindaklanjutDomain.RiwayatPemeriksaanJoinRow
	for pgxRows.Next() {
		var row tindaklanjutDomain.RiwayatPemeriksaanJoinRow
		err := pgxRows.Scan(&row.Tanggal, &row.BeratBadan, &row.TinggiBadan)
		if err != nil {
			return nil, err
		}
		rows = append(rows, &row)
	}

	return rows, nil
}

func (r *Repo) GetTindakLanjutByHasilPemeriksaan(ctx context.Context, idHasilPemeriksaan int32) (*model.TindakLanjut, error) {
	var results []*model.TindakLanjut
	stmt := SELECT(TindakLanjut.AllColumns).
		FROM(TindakLanjut).
		WHERE(TindakLanjut.IDHasilPemeriksaan.EQ(Int32(idHasilPemeriksaan)))

	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}

func (r *Repo) CreateTindakLanjut(ctx context.Context, data *model.TindakLanjut) error {
	stmt := TindakLanjut.INSERT(
		TindakLanjut.IDHasilPemeriksaan,
		TindakLanjut.IDBidan,
		TindakLanjut.CatatanMedis,
		TindakLanjut.Rekomendasi,
		TindakLanjut.JadwalKontrol,
		TindakLanjut.StatusPasien,
		TindakLanjut.CreatedAt,
		TindakLanjut.UpdatedAt,
	).MODEL(data).
		RETURNING(TindakLanjut.IDTindakLanjut)

	var results []*model.TindakLanjut
	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return err
	}
	if len(results) > 0 {
		data.IDTindakLanjut = results[0].IDTindakLanjut
	}
	return nil
}

func (r *Repo) CreateRujukan(ctx context.Context, data *model.Rujukan) error {
	stmt := Rujukan.INSERT(
		Rujukan.IDTindakLanjut,
		Rujukan.AlasanRujukan,
		Rujukan.TanggalRujukan,
		Rujukan.StatusRujukan,
		Rujukan.IDFaskes,
		Rujukan.CreatedAt,
		Rujukan.UpdatedAt,
	).MODEL(data).
		RETURNING(Rujukan.IDRujukan)

	var results []*model.Rujukan
	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return err
	}
	if len(results) > 0 {
		data.IDRujukan = results[0].IDRujukan
	}
	return nil
}

func (r *Repo) UpdateStatusRujukan(ctx context.Context, idRujukan int32, statusRujukan model.StatusRujukan) (*model.Rujukan, error) {
	var results []*model.Rujukan
	stmt := Rujukan.UPDATE(
		Rujukan.StatusRujukan,
		Rujukan.UpdatedAt,
	).SET(
		statusRujukanExpr[statusRujukan],
		TimestampzT(time.Now()),
	).WHERE(Rujukan.IDRujukan.EQ(Int32(idRujukan))).
		RETURNING(Rujukan.AllColumns)

	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}

func (r *Repo) GetStatusTindakLanjut(ctx context.Context) ([]*tindaklanjutDomain.StatusTindakLanjutJoinRow, error) {
	sql := `
		SELECT
			COALESCE(p.id_pasien, 0) AS id_pasien,
			COALESCE(ua.nama, '') AS nama_pasien,
			tl.status_pasien::text AS status_pasien,
			COALESCE(r.status_rujukan::text, '') AS status_rujukan,
			COALESCE(r.tanggal_rujukan::text, '') AS tanggal_rujukan
		FROM tindak_lanjut tl
		LEFT JOIN hasil_pemeriksaan hp ON hp.id_hasil_pemeriksaan = tl.id_hasil_pemeriksaan
		LEFT JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
		LEFT JOIN pasien p ON p.id_pasien = ji.id_pasien AND p.is_deleted = false
		LEFT JOIN user_account ua ON ua.id_user = p.id_pasien AND ua.is_deleted = false
		LEFT JOIN rujukan r ON r.id_tindak_lanjut = tl.id_tindak_lanjut
		ORDER BY tl.created_at DESC
	`

	pgxRows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("GetStatusTindakLanjut query failed: %w", err)
	}
	defer pgxRows.Close()

	var rows []*tindaklanjutDomain.StatusTindakLanjutJoinRow
	for pgxRows.Next() {
		var row tindaklanjutDomain.StatusTindakLanjutJoinRow
		err := pgxRows.Scan(&row.IDPasien, &row.NamaPasien, &row.StatusPasien, &row.StatusRujukan, &row.TanggalRujukan)
		if err != nil {
			return nil, err
		}
		rows = append(rows, &row)
	}

	return rows, nil
}

func (r *Repo) GetStatusTindakLanjutByUserID(ctx context.Context, idUser int32) ([]*tindaklanjutDomain.StatusTindakLanjutJoinRow, error) {
	sql := `
		SELECT
			p.id_pasien,
			COALESCE(a.nama_anak, ua.nama) AS nama_pasien,
			tl.status_pasien::text,
			COALESCE(r.status_rujukan::text, '') AS status_rujukan,
			COALESCE(r.tanggal_rujukan::text, '') AS tanggal_rujukan
		FROM tindak_lanjut tl
		JOIN hasil_pemeriksaan hp ON hp.id_hasil_pemeriksaan = tl.id_hasil_pemeriksaan
		JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
		JOIN pasien p ON p.id_pasien = ji.id_pasien AND p.is_deleted = false
		JOIN user_account ua ON ua.id_user = p.id_pasien AND ua.is_deleted = false
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
		LEFT JOIN rujukan r ON r.id_tindak_lanjut = tl.id_tindak_lanjut
		WHERE (a.id_wali = #1 OR p.id_pasien = #1)
		ORDER BY tl.created_at DESC
	`

	var rows []*tindaklanjutDomain.StatusTindakLanjutJoinRow
	err := pgxV5.Query(ctx, RawStatement(sql, RawArgs{"#1": idUser}), r.db, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repo) GetLaporanTindakLanjut(ctx context.Context) ([]*tindaklanjutDomain.LaporanTindakLanjutJoinRow, error) {
	sql := `
		SELECT
			COALESCE(l.nama_lokasi, 'Unknown') AS wilayah,
			COUNT(DISTINCT r.id_rujukan) FILTER (WHERE r.status_rujukan = 'Diajukan') AS jumlah_pasien_dirujuk,
			COUNT(DISTINCT r.id_rujukan) FILTER (WHERE r.status_rujukan = 'Diterima') AS jumlah_pasien_diterima,
			COUNT(DISTINCT r.id_rujukan) FILTER (WHERE r.status_rujukan = 'Diproses') AS jumlah_pasien_diproses
		FROM rujukan r
		JOIN tindak_lanjut tl ON tl.id_tindak_lanjut = r.id_tindak_lanjut
		JOIN hasil_pemeriksaan hp ON hp.id_hasil_pemeriksaan = tl.id_hasil_pemeriksaan
		JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
		JOIN pasien p ON p.id_pasien = ji.id_pasien AND p.is_deleted = false
		LEFT JOIN posyandu pos ON pos.id_posyandu = p.id_posyandu
		LEFT JOIN lokasi l ON l.id_lokasi = pos.id_lokasi
		GROUP BY l.nama_lokasi
		ORDER BY l.nama_lokasi
	`

	pgxRows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer pgxRows.Close()

	var rows []*tindaklanjutDomain.LaporanTindakLanjutJoinRow
	for pgxRows.Next() {
		var row tindaklanjutDomain.LaporanTindakLanjutJoinRow
		err := pgxRows.Scan(&row.Wilayah, &row.JumlahPasienDirujuk, &row.JumlahPasienDiterima, &row.JumlahPasienDiproses)
		if err != nil {
			return nil, err
		}
		rows = append(rows, &row)
	}

	return rows, nil
}

func (r *Repo) GetDetailTindakLanjutByID(ctx context.Context, idTindakLanjut int32) (*tindaklanjutDomain.DetailTindakLanjutJoinRow, error) {
	sql := `
		SELECT
			tl.id_tindak_lanjut,
			tl.status_pasien::text AS status_pasien,
			tl.catatan_medis,
			tl.rekomendasi,
			tl.jadwal_kontrol::text AS jadwal_kontrol,
			r.status_rujukan::text AS status_rujukan,
			fk.nama_faskes
		FROM tindak_lanjut tl
		LEFT JOIN rujukan r ON r.id_tindak_lanjut = tl.id_tindak_lanjut
		LEFT JOIN fasilitas_kesehatan fk ON fk.id_faskes = r.id_faskes
		WHERE tl.id_tindak_lanjut = $1
		LIMIT 1
	`

	pgxRows, err := r.db.Query(ctx, sql, idTindakLanjut)
	if err != nil {
		return nil, err
	}
	defer pgxRows.Close()

	if !pgxRows.Next() {
		return nil, nil
	}

	var row tindaklanjutDomain.DetailTindakLanjutJoinRow
	err = pgxRows.Scan(&row.IDTindakLanjut, &row.StatusPasien, &row.CatatanMedis, &row.Rekomendasi, &row.JadwalKontrol, &row.StatusRujukan, &row.NamaFaskes)
	if err != nil {
		return nil, err
	}
	return &row, nil
}
