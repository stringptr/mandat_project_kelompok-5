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

func (r *Repo) GetAll(ctx context.Context) ([]*imunisasiDomain.ImunisasiJoinRow, error) {
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
		ORDER BY ji.tanggal_jadwal DESC
	`

	var rows []*imunisasiDomain.ImunisasiJoinRow
	err := pgxV5.Query(ctx, RawStatement(sql), r.db, &rows)
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
			ji.tanggal_jadwal::text,
			ji.tanggal_realisasi::text,
			ji.status_imunisasi::text
		FROM jadwal_imunisasi ji
		JOIN pasien p ON p.id_pasien = ji.id_pasien
		JOIN user_account ua ON ua.id_user = p.id_pasien
		LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
		WHERE ji.id_imunisasi = #1
		LIMIT 1
	`

	var row imunisasiDomain.DetailJoinRow
	err := pgxV5.Query(ctx, RawStatement(sql, RawArgs{"#1": idImunisasi}), r.db, &row)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, nil
		}
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

	var row imunisasiDomain.StatistikRow
	err := pgxV5.Query(ctx, RawStatement(sql), r.db, &row)
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
