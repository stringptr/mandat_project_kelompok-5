package pemeriksaan

import (
	"context"
	"fmt"
	"strings"

	pemeriksaanDomain "github.com/stringptr/SiGizi/backend/internal/domain/pemeriksaan"
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

func (r *Repo) GetByID(ctx context.Context, idHasilPemeriksaan int32) (*model.HasilPemeriksaan, error) {
	var results []*model.HasilPemeriksaan
	stmt := SELECT(HasilPemeriksaan.AllColumns).
		FROM(HasilPemeriksaan).
		WHERE(HasilPemeriksaan.IDHasilPemeriksaan.EQ(Int32(idHasilPemeriksaan)))
	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}

func (r *Repo) GetJadwalImunisasiByID(ctx context.Context, idJadwalImunisasi int32) (*model.JadwalImunisasi, error) {
	var results []*model.JadwalImunisasi
	stmt := SELECT(JadwalImunisasi.AllColumns).
		FROM(JadwalImunisasi).
		WHERE(JadwalImunisasi.IDImunisasi.EQ(Int32(idJadwalImunisasi)))
	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}

func (r *Repo) GetUserByID(ctx context.Context, idUser int32) (*model.UserAccount, error) {
	var results []*model.UserAccount
	stmt := SELECT(UserAccount.AllColumns).
		FROM(UserAccount).
		WHERE(UserAccount.IDUser.EQ(Int32(idUser)).AND(UserAccount.IsDeleted.EQ(Bool(false))))
	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0], nil
}

func (r *Repo) GetNamaPetugasByID(ctx context.Context, idUser int32) (string, error) {
	var results []struct {
		Nama string
	}
	stmt := SELECT(UserAccount.Nama).
		FROM(UserAccount).
		WHERE(UserAccount.IDUser.EQ(Int32(idUser)))
	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", nil
	}
	return results[0].Nama, nil
}

func (r *Repo) Create(ctx context.Context, data *model.HasilPemeriksaan) error {
	stmt := HasilPemeriksaan.INSERT(
		HasilPemeriksaan.IDPetugasInput,
		HasilPemeriksaan.IDJadwalImunisasi,
		HasilPemeriksaan.BeratBadan,
		HasilPemeriksaan.TinggiBadan,
		HasilPemeriksaan.LingkarKepala,
		HasilPemeriksaan.TekananDarah,
		HasilPemeriksaan.StatusStunting,
		HasilPemeriksaan.StatusGizi,
		HasilPemeriksaan.Catatan,
		HasilPemeriksaan.CreatedAt,
		HasilPemeriksaan.UpdatedAt,
	).MODEL(data).
		RETURNING(HasilPemeriksaan.IDHasilPemeriksaan)

	var results []*model.HasilPemeriksaan
	err := pgxV5.Query(ctx, stmt, r.db, &results)
	if err != nil {
		return err
	}
	if len(results) > 0 {
		data.IDHasilPemeriksaan = results[0].IDHasilPemeriksaan
	}
	return nil
}

func (r *Repo) Update(ctx context.Context, data *model.HasilPemeriksaan) error {
	stmt := HasilPemeriksaan.UPDATE(
		HasilPemeriksaan.BeratBadan,
		HasilPemeriksaan.TinggiBadan,
		HasilPemeriksaan.LingkarKepala,
		HasilPemeriksaan.TekananDarah,
		HasilPemeriksaan.StatusStunting,
		HasilPemeriksaan.StatusGizi,
		HasilPemeriksaan.Catatan,
		HasilPemeriksaan.UpdatedAt,
	).MODEL(data).
		WHERE(HasilPemeriksaan.IDHasilPemeriksaan.EQ(Int32(data.IDHasilPemeriksaan)))

	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) Delete(ctx context.Context, idHasilPemeriksaan int32) error {
	stmt := HasilPemeriksaan.DELETE().
		WHERE(HasilPemeriksaan.IDHasilPemeriksaan.EQ(Int32(idHasilPemeriksaan)))

	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) GetAllPaginated(ctx context.Context, page, perPage int, q string) ([]*pemeriksaanDomain.GetAllJoinRow, int, error) {
	offset := (page - 1) * perPage

	countSQL := `
		SELECT COUNT(*)
		FROM hasil_pemeriksaan hp
		JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
		JOIN pasien p ON p.id_pasien = ji.id_pasien
		JOIN user_account ua ON ua.id_user = p.id_pasien
		JOIN user_account petugas ON petugas.id_user = hp.id_petugas_input
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
	`

	dataSQL := `
		SELECT
			hp.id_hasil_pemeriksaan,
			COALESCE(a.nama_anak, ua.nama) AS nama_pasien,
			petugas.nama AS diinput_oleh,
			hp.status_stunting::text,
			hp.status_gizi::text,
			hp.created_at AS tanggal_input
		FROM hasil_pemeriksaan hp
		JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
		JOIN pasien p ON p.id_pasien = ji.id_pasien
		JOIN user_account ua ON ua.id_user = p.id_pasien
		JOIN user_account petugas ON petugas.id_user = hp.id_petugas_input
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
	`

	args := map[string]any{}
	if q != "" {
		filter := ` WHERE (ua.nama ILIKE '%' || #q || '%' OR COALESCE(a.nama_anak, '') ILIKE '%' || #q || '%')`
		countSQL += filter
		dataSQL += filter
		args["#q"] = q
	}

	var total int
	err := pgxV5.Query(ctx, RawStatement(countSQL, RawArgs(args)), r.db, &total)
	if err != nil {
		return nil, 0, err
	}

	dataSQL += ` ORDER BY hp.created_at DESC`
	dataSQL += fmt.Sprintf(` LIMIT %d OFFSET %d`, perPage, offset)

	var rows []*pemeriksaanDomain.GetAllJoinRow
	err = pgxV5.Query(ctx, RawStatement(dataSQL, RawArgs(args)), r.db, &rows)
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *Repo) CheckPemeriksaanOwnership(ctx context.Context, idHasilPemeriksaan int32, idUser int32) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM hasil_pemeriksaan hp
			JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
			JOIN pasien p ON p.id_pasien = ji.id_pasien
			LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
			WHERE hp.id_hasil_pemeriksaan = #1
			  AND p.is_deleted = false
			  AND (p.id_pasien = #2 OR a.id_wali = #2)
		) AS owned
	`
	var result struct{ Owned bool }
	err := pgxV5.Query(ctx, RawStatement(query, RawArgs{"#1": idHasilPemeriksaan, "#2": idUser}), r.db, &result)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return false, nil
		}
		return false, err
	}
	return result.Owned, nil
}

func (r *Repo) GetPendingVerification(ctx context.Context, page int, perPage int) ([]*pemeriksaanDomain.PendingJoinRow, int, error) {
	offset := (page - 1) * perPage

	selectFrom := `
		FROM hasil_pemeriksaan hp
		JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
		JOIN pasien p ON p.id_pasien = ji.id_pasien AND p.is_deleted = false
		JOIN user_account ua ON ua.id_user = p.id_pasien
		JOIN user_account petugas ON petugas.id_user = hp.id_petugas_input
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
	`

	var countResult struct{ Count int64 }
	err := pgxV5.Query(ctx, RawStatement("SELECT COUNT(*)"+selectFrom, RawArgs{}), r.db, &countResult)
	if err != nil {
		return nil, 0, err
	}

	dataSQL := `
		SELECT
			hp.id_hasil_pemeriksaan,
			COALESCE(a.nama_anak, ua.nama) AS nama_pasien,
			petugas.nama AS diinput_oleh,
			hp.created_at AS tanggal_input
	` + selectFrom + `
		ORDER BY hp.created_at DESC
		LIMIT $1 OFFSET $2
	`

	pgxRows, err := r.db.Query(ctx, dataSQL, perPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer pgxRows.Close()

	var rows []*pemeriksaanDomain.PendingJoinRow
	for pgxRows.Next() {
		var row pemeriksaanDomain.PendingJoinRow
		err := pgxRows.Scan(&row.IDHasilPemeriksaan, &row.NamaPasien, &row.DiinputOleh, &row.TanggalInput)
		if err != nil {
			return nil, 0, err
		}
		rows = append(rows, &row)
	}

	return rows, int(countResult.Count), nil
}

func (r *Repo) GetDetailJoinByID(ctx context.Context, idHasilPemeriksaan int32) (*pemeriksaanDomain.DetailJoinRow, error) {
	sql := `
		SELECT
			hp.id_hasil_pemeriksaan,
			p.id_pasien,
			COALESCE(a.nama_anak, ua.nama) AS nama_pasien,
			hp.berat_badan::float8 AS berat_badan,
			hp.tinggi_badan::float8 AS tinggi_badan,
			hp.lingkar_kepala::float8 AS lingkar_kepala,
			hp.tekanan_darah,
			hp.status_stunting::text AS status_stunting,
			hp.status_gizi::text AS status_gizi,
			hp.catatan
		FROM hasil_pemeriksaan hp
		JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
		JOIN pasien p ON p.id_pasien = ji.id_pasien
		JOIN user_account ua ON ua.id_user = p.id_pasien
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
		WHERE hp.id_hasil_pemeriksaan = $1
		LIMIT 1
	`

	pgxRows, err := r.db.Query(ctx, sql, idHasilPemeriksaan)
	if err != nil {
		return nil, err
	}
	defer pgxRows.Close()

	if !pgxRows.Next() {
		return nil, nil
	}

	var row pemeriksaanDomain.DetailJoinRow
	err = pgxRows.Scan(&row.IDHasilPemeriksaan, &row.IDPasien, &row.NamaPasien, &row.BeratBadan, &row.TinggiBadan, &row.LingkarKepala, &row.TekananDarah, &row.StatusStunting, &row.StatusGizi, &row.Catatan)
	if err != nil {
		return nil, err
	}
	return &row, nil
}


