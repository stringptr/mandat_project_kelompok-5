package pemeriksaan

import (
	"context"
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

func (r *Repo) GetPendingVerification(ctx context.Context) ([]*pemeriksaanDomain.PendingJoinRow, error) {
	sql := `
		SELECT
			hp.id_hasil_pemeriksaan,
			COALESCE(a.nama_anak, ua.nama) AS nama_pasien,
			petugas.nama AS diinput_oleh,
			hp.created_at AS tanggal_input
		FROM hasil_pemeriksaan hp
		JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
		JOIN pasien p ON p.id_pasien = ji.id_pasien
		JOIN user_account ua ON ua.id_user = p.id_pasien
		JOIN user_account petugas ON petugas.id_user = hp.id_petugas_input
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
		ORDER BY hp.created_at DESC
	`

	var rows []*pemeriksaanDomain.PendingJoinRow
	err := pgxV5.Query(ctx, RawStatement(sql), r.db, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repo) GetDetailJoinByID(ctx context.Context, idHasilPemeriksaan int32) (*pemeriksaanDomain.DetailJoinRow, error) {
	sql := `
		SELECT
			hp.id_hasil_pemeriksaan,
			p.id_pasien,
			COALESCE(a.nama_anak, ua.nama) AS nama_pasien,
			hp.berat_badan::float8,
			hp.tinggi_badan::float8,
			hp.lingkar_kepala::float8,
			hp.tekanan_darah,
			hp.status_stunting::text,
			hp.status_gizi::text,
			hp.catatan
		FROM hasil_pemeriksaan hp
		JOIN jadwal_imunisasi ji ON ji.id_imunisasi = hp.id_jadwal_imunisasi
		JOIN pasien p ON p.id_pasien = ji.id_pasien
		JOIN user_account ua ON ua.id_user = p.id_pasien
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
		WHERE hp.id_hasil_pemeriksaan = #1
		LIMIT 1
	`

	var row pemeriksaanDomain.DetailJoinRow
	err := pgxV5.Query(ctx, RawStatement(sql, RawArgs{"#1": idHasilPemeriksaan}), r.db, &row)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}


