package imunisasi

import (
	"context"
	"strings"
	"time"

	imunisasiDomain "github.com/stringptr/SiGizi/backend/internal/domain/imunisasi"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/enum"
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

func (r *Repo) GetAllByUser(ctx context.Context, idUser int32, page int, perPage int, q string) ([]*imunisasiDomain.ImunisasiJoinRow, int, error) {
	offset := int64((page - 1) * perPage)

	selectFrom := `
		FROM jadwal_imunisasi ji
		JOIN pasien p ON p.id_pasien = ji.id_pasien AND p.is_deleted = false
		JOIN user_account ua ON ua.id_user = p.id_pasien
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
		WHERE (p.id_pasien = #1 OR a.id_wali = #1)
	`

	var countResult struct{ Count int64 }
	err := pgxV5.Query(ctx, RawStatement("SELECT COUNT(*)"+selectFrom, RawArgs{"#1": idUser}), r.db, &countResult)
	if err != nil {
		return nil, 0, err
	}

	dataSQL := `
		SELECT
			ji.id_imunisasi,
			COALESCE(a.nama_anak, ua.nama) AS nama_pasien,
			ji.nama_vaksin,
			ji.tanggal_jadwal::text AS tanggal_jadwal,
			ji.status_imunisasi::text AS status_imunisasi
	` + strings.ReplaceAll(selectFrom, "#1", "$1") + `
		ORDER BY ji.tanggal_jadwal DESC
		OFFSET $2 LIMIT $3
	`

	pgxRows, err := r.db.Query(ctx, dataSQL, idUser, offset, perPage)
	if err != nil {
		return nil, 0, err
	}
	defer pgxRows.Close()

	var rows []*imunisasiDomain.ImunisasiJoinRow
	for pgxRows.Next() {
		var row imunisasiDomain.ImunisasiJoinRow
		err := pgxRows.Scan(&row.IDImunisasi, &row.NamaPasien, &row.NamaVaksin, &row.TanggalJadwal, &row.StatusImunisasi)
		if err != nil {
			return nil, 0, err
		}
		rows = append(rows, &row)
	}

	return rows, int(countResult.Count), nil
}

func (r *Repo) GetAll(ctx context.Context, page int, perPage int, q string) ([]*imunisasiDomain.ImunisasiJoinRow, int, error) {
	offset := int64((page - 1) * perPage)

	selectFrom := `
		FROM jadwal_imunisasi ji
		JOIN pasien p ON p.id_pasien = ji.id_pasien AND p.is_deleted = false
		JOIN user_account ua ON ua.id_user = p.id_pasien
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
	`

	var countResult struct{ Count int64 }
	err := pgxV5.Query(ctx, RawStatement("SELECT COUNT(*)"+selectFrom), r.db, &countResult)
	if err != nil {
		return nil, 0, err
	}

	dataSQL := `
		SELECT
			ji.id_imunisasi,
			COALESCE(a.nama_anak, ua.nama) AS nama_pasien,
			ji.nama_vaksin,
			ji.tanggal_jadwal::text AS tanggal_jadwal,
			ji.status_imunisasi::text AS status_imunisasi
	` + selectFrom + `
		ORDER BY ji.tanggal_jadwal DESC
		OFFSET $1 LIMIT $2
	`

	pgxRows, err := r.db.Query(ctx, dataSQL, offset, perPage)
	if err != nil {
		return nil, 0, err
	}
	defer pgxRows.Close()

	var rows []*imunisasiDomain.ImunisasiJoinRow
	for pgxRows.Next() {
		var row imunisasiDomain.ImunisasiJoinRow
		err := pgxRows.Scan(&row.IDImunisasi, &row.NamaPasien, &row.NamaVaksin, &row.TanggalJadwal, &row.StatusImunisasi)
		if err != nil {
			return nil, 0, err
		}
		rows = append(rows, &row)
	}

	return rows, int(countResult.Count), nil
}

func (r *Repo) GetAllByUserID(ctx context.Context, idUser int32) ([]*imunisasiDomain.ImunisasiJoinRow, error) {
	sql := `
		SELECT
			ji.id_imunisasi,
			COALESCE(a.nama_anak, ua.nama) AS nama_pasien,
			ji.nama_vaksin,
			ji.tanggal_jadwal::text,
			ji.status_imunisasi::text
		FROM jadwal_imunisasi ji
		JOIN pasien p ON p.id_pasien = ji.id_pasien
		JOIN user_account ua ON ua.id_user = p.id_pasien
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
		WHERE (a.id_wali = #1 OR p.id_pasien = #1)
		ORDER BY ji.tanggal_jadwal DESC
	`

	var rows []*imunisasiDomain.ImunisasiJoinRow
	err := pgxV5.Query(ctx, RawStatement(sql, RawArgs{"#1": idUser}), r.db, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repo) GetByID(ctx context.Context, idImunisasi int32) (*model.JadwalImunisasi, error) {
	var results []*model.JadwalImunisasi
	stmt := SELECT(JadwalImunisasi.AllColumns).
		FROM(JadwalImunisasi).
		WHERE(JadwalImunisasi.IDImunisasi.EQ(Int32(idImunisasi)))
	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}

func (r *Repo) GetDetailJoinByID(ctx context.Context, idImunisasi int32) (*imunisasiDomain.DetailJoinRow, error) {
	sql := `
		SELECT
			ji.id_imunisasi,
			p.id_pasien,
			COALESCE(a.nama_anak, ua.nama) AS nama_pasien,
			ji.nama_vaksin,
			ji.tanggal_jadwal::text AS tanggal_jadwal,
			ji.tanggal_realisasi::text AS tanggal_realisasi,
			ji.status_imunisasi::text AS status_imunisasi
		FROM jadwal_imunisasi ji
		JOIN pasien p ON p.id_pasien = ji.id_pasien
		JOIN user_account ua ON ua.id_user = p.id_pasien
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
		WHERE ji.id_imunisasi = $1
		LIMIT 1
	`

	pgxRows, err := r.db.Query(ctx, sql, idImunisasi)
	if err != nil {
		return nil, err
	}
	defer pgxRows.Close()

	if !pgxRows.Next() {
		return nil, nil
	}

	var row imunisasiDomain.DetailJoinRow
	err = pgxRows.Scan(&row.IDImunisasi, &row.IDPasien, &row.NamaPasien, &row.NamaVaksin, &row.TanggalJadwal, &row.TanggalRealisasi, &row.StatusImunisasi)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) GetByPasienID(ctx context.Context, idPasien int32) ([]*model.JadwalImunisasi, error) {
	var results []*model.JadwalImunisasi
	stmt := SELECT(JadwalImunisasi.AllColumns).
		FROM(JadwalImunisasi).
		WHERE(JadwalImunisasi.IDPasien.EQ(Int32(idPasien))).
		ORDER_BY(JadwalImunisasi.TanggalJadwal.ASC())
	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

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

func (r *Repo) CheckPasienOwnership(ctx context.Context, idPasien int32, idUser int32) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM pasien p
			LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
			WHERE p.id_pasien = #1
			  AND p.is_deleted = false
			  AND (p.id_pasien = #2 OR a.id_wali = #2)
		) AS owned
	`
	var result struct{ Owned bool }
	err := pgxV5.Query(ctx, RawStatement(query, RawArgs{"#1": idPasien, "#2": idUser}), r.db, &result)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return false, nil
		}
		return false, err
	}
	return result.Owned, nil
}

func (r *Repo) GetAnakByPasienID(ctx context.Context, idPasien int32) (*model.Anak, error) {
	var results []*model.Anak
	stmt := SELECT(Anak.AllColumns).
		FROM(Anak).
		WHERE(Anak.IDPasien.EQ(Int32(idPasien)).AND(Anak.IsDeleted.EQ(Bool(false))))
	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}

func (r *Repo) GetNamaPasienByID(ctx context.Context, idPasien int32) (string, error) {
	var results []struct {
		Nama string
	}
	stmt := SELECT(UserAccount.Nama).
		FROM(UserAccount).
		WHERE(UserAccount.IDUser.EQ(Int32(idPasien)))
	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", nil
	}
	return results[0].Nama, nil
}

func (r *Repo) Create(ctx context.Context, data *model.JadwalImunisasi) error {
	stmt := JadwalImunisasi.INSERT(
		JadwalImunisasi.IDPasien,
		JadwalImunisasi.NamaVaksin,
		JadwalImunisasi.TanggalJadwal,
		JadwalImunisasi.StatusImunisasi,
		JadwalImunisasi.CreatedAt,
		JadwalImunisasi.UpdatedAt,
	).MODEL(data).
		RETURNING(JadwalImunisasi.IDImunisasi)

	var results []*model.JadwalImunisasi
	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return err
	}
	if len(results) > 0 {
		data.IDImunisasi = results[0].IDImunisasi
	}
	return nil
}

func (r *Repo) Update(ctx context.Context, data *model.JadwalImunisasi) error {
	stmt := JadwalImunisasi.UPDATE(
		JadwalImunisasi.IDPasien,
		JadwalImunisasi.NamaVaksin,
		JadwalImunisasi.TanggalJadwal,
		JadwalImunisasi.UpdatedAt,
	).MODEL(data).
		WHERE(JadwalImunisasi.IDImunisasi.EQ(Int32(data.IDImunisasi)))

	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) Delete(ctx context.Context, idImunisasi int32) error {
	stmt := JadwalImunisasi.DELETE().
		WHERE(JadwalImunisasi.IDImunisasi.EQ(Int32(idImunisasi)))

	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) UpdateRealisasi(ctx context.Context, idImunisasi int32, tanggalRealisasi string) error {
	t, err := time.Parse("2006-01-02", tanggalRealisasi)
	if err != nil {
		return err
	}
	stmt := JadwalImunisasi.UPDATE(
		JadwalImunisasi.TanggalRealisasi,
		JadwalImunisasi.StatusImunisasi,
		JadwalImunisasi.UpdatedAt,
	).SET(
		TimestampzT(t),
		enum.StatusImunisasi.Sudah,
		RawTimestampz("NOW()"),
	).WHERE(JadwalImunisasi.IDImunisasi.EQ(Int32(idImunisasi)))

	_, err = pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) GetStatistik(ctx context.Context) (*imunisasiDomain.StatistikRow, error) {
	sql := `
		SELECT
			COUNT(*) AS total_target,
			COUNT(*) FILTER (WHERE status_imunisasi = 'Sudah') AS total_realisasi
		FROM jadwal_imunisasi
	`

	pgxRows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer pgxRows.Close()

	if !pgxRows.Next() {
		return nil, nil
	}

	var row imunisasiDomain.StatistikRow
	err = pgxRows.Scan(&row.TotalTarget, &row.TotalRealisasi)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *Repo) GetVaksinTerbanyak(ctx context.Context) (string, error) {
	var results []struct {
		NamaVaksin string
	}
	stmt := SELECT(
		JadwalImunisasi.NamaVaksin,
	).FROM(
		JadwalImunisasi,
	).GROUP_BY(
		JadwalImunisasi.NamaVaksin,
	).ORDER_BY(
		COUNT(STAR).DESC(),
	).LIMIT(1)

	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", nil
	}
	return results[0].NamaVaksin, nil
}
