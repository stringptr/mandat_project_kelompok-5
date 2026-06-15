package tindaklanjut

import (
	"context"
	"strings"
	"time"

	tindaklanjutDomain "github.com/stringptr/SiGizi/backend/internal/domain/tindaklanjut"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
	. "github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/table"

	"github.com/go-jet/jet/v2/pgxV5"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetPasienTindakLanjut(ctx context.Context) ([]*tindaklanjutDomain.PasienTindakLanjutJoinRow, error) {
	sql := `
		SELECT
			p.id_pasien,
			COALESCE(a.nama_anak, ua.nama) AS nama_pasien,
			hp.status_gizi::text,
			CASE
				WHEN hp.status_gizi IN ('Gizi Buruk', 'Gizi Kurang')
					OR hp.status_stunting IN ('Stunting', 'Stunting Berat', 'Berisiko Stunting')
				THEN 'Perlu Rujukan'
				ELSE 'Dalam Pemantauan'
			END AS status_pasien,
			hp.created_at::text AS tanggal_pemeriksaan
		FROM pasien p
		JOIN user_account ua ON ua.id_user = p.id_pasien AND ua.is_deleted = false
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
		LEFT JOIN hasil_pemeriksaan hp ON hp.id_hasil_pemeriksaan = (
			SELECT hp2.id_hasil_pemeriksaan
			FROM hasil_pemeriksaan hp2
			JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp2.id_jadwal_imunisasi
			WHERE ji.id_pasien = p.id_pasien
			ORDER BY hp2.created_at DESC
			LIMIT 1
		)
		LEFT JOIN tindak_lanjut tl ON tl.id_hasil_pemeriksaan = hp.id_hasil_pemeriksaan
		WHERE p.is_deleted = false
			AND hp.status_gizi IS NOT NULL
			AND hp.status_gizi != 'Gizi Baik'
			AND tl.id_tindak_lanjut IS NULL
		ORDER BY hp.created_at DESC
	`

	var rows []*tindaklanjutDomain.PasienTindakLanjutJoinRow
	err := pgxV5.Query(ctx, RawStatement(sql), r.db, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repo) GetDetailPasienByID(ctx context.Context, idPasien int32) (*tindaklanjutDomain.DetailPasienJoinRow, error) {
	sql := `
		SELECT
			p.id_pasien,
			COALESCE(a.nama_anak, ua.nama) AS nama_pasien,
			EXTRACT(YEAR FROM AGE(ua.tanggal_lahir))::int || ' Tahun' AS usia,
			hp.status_gizi::text,
			hp.status_stunting::text,
			COALESCE(hp.catatan, '') AS catatan
		FROM pasien p
		JOIN user_account ua ON ua.id_user = p.id_pasien AND ua.is_deleted = false
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
		LEFT JOIN hasil_pemeriksaan hp ON hp.id_hasil_pemeriksaan = (
			SELECT hp2.id_hasil_pemeriksaan
			FROM hasil_pemeriksaan hp2
			JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp2.id_jadwal_imunisasi
			WHERE ji.id_pasien = p.id_pasien
			ORDER BY hp2.created_at DESC
			LIMIT 1
		)
		WHERE p.id_pasien = #1 AND p.is_deleted = false
		LIMIT 1
	`

	var row tindaklanjutDomain.DetailPasienJoinRow
	err := pgxV5.Query(ctx, RawStatement(sql, RawArgs{"#1": idPasien}), r.db, &row)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

func (r *Repo) GetRiwayatPemeriksaan(ctx context.Context, idPasien int32) ([]*tindaklanjutDomain.RiwayatPemeriksaanJoinRow, error) {
	sql := `
		SELECT
			hp.created_at::text AS tanggal,
			hp.berat_badan::float8,
			hp.tinggi_badan::float8
		FROM hasil_pemeriksaan hp
		JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
		WHERE ji.id_pasien = #1
		ORDER BY hp.created_at DESC
	`

	var rows []*tindaklanjutDomain.RiwayatPemeriksaanJoinRow
	err := pgxV5.Query(ctx, RawStatement(sql, RawArgs{"#1": idPasien}), r.db, &rows)
	if err != nil {
		return nil, err
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
		String(string(statusRujukan)),
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
		ORDER BY tl.created_at DESC
	`

	var rows []*tindaklanjutDomain.StatusTindakLanjutJoinRow
	err := pgxV5.Query(ctx, RawStatement(sql), r.db, &rows)
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
		JOIN posyandu pos ON pos.id_posyandu = p.id_posyandu
		JOIN lokasi l ON l.id_lokasi = pos.id_lokasi
		GROUP BY l.nama_lokasi
		ORDER BY l.nama_lokasi
	`

	var rows []*tindaklanjutDomain.LaporanTindakLanjutJoinRow
	err := pgxV5.Query(ctx, RawStatement(sql), r.db, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repo) GetDetailTindakLanjutByID(ctx context.Context, idTindakLanjut int32) (*tindaklanjutDomain.DetailTindakLanjutJoinRow, error) {
	sql := `
		SELECT
			tl.id_tindak_lanjut,
			tl.status_pasien::text,
			tl.catatan_medis,
			tl.rekomendasi,
			tl.jadwal_kontrol::text,
			r.status_rujukan::text,
			fk.nama_faskes
		FROM tindak_lanjut tl
		LEFT JOIN rujukan r ON r.id_tindak_lanjut = tl.id_tindak_lanjut
		LEFT JOIN fasilitas_kesehatan fk ON fk.id_faskes = r.id_faskes
		WHERE tl.id_tindak_lanjut = #1
		LIMIT 1
	`

	var row tindaklanjutDomain.DetailTindakLanjutJoinRow
	err := pgxV5.Query(ctx, RawStatement(sql, RawArgs{"#1": idTindakLanjut}), r.db, &row)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}
