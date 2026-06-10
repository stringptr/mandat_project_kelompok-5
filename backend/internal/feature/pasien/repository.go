package pasien

import (
	"context"
	"fmt"
	"strings"
	"time"

	pasienDomain "github.com/stringptr/SiGizi/backend/internal/domain/pasien"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
	. "github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/table"

	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/pgxV5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetUserByID(ctx context.Context, idUser int32) (*model.UserAccount, error) {
	var users []*model.UserAccount
	stmt := SELECT(UserAccount.AllColumns).FROM(UserAccount).WHERE(UserAccount.IDUser.EQ(Int32(idUser)))
	err := pgxV5.Query(ctx, stmt, r.db, &users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

func (r *Repo) GetPasienByUserID(ctx context.Context, idUser int32) (*model.Pasien, error) {
	var pasien []*model.Pasien
	stmt := SELECT(Pasien.AllColumns).FROM(Pasien).WHERE(Pasien.IDPasien.EQ(Int32(idUser)))
	err := pgxV5.Query(ctx, stmt, r.db, &pasien)
	if err != nil {
		return nil, err
	}
	if len(pasien) == 0 {
		return nil, nil
	}
	return pasien[0], nil
}

func (r *Repo) GetPosyanduByID(ctx context.Context, idPosyandu int32) (*model.Posyandu, error) {
	var posyandu []*model.Posyandu
	stmt := SELECT(Posyandu.AllColumns).FROM(Posyandu).WHERE(Posyandu.IDPosyandu.EQ(Int32(idPosyandu)))
	err := pgxV5.Query(ctx, stmt, r.db, &posyandu)
	if err != nil {
		return nil, err
	}
	if len(posyandu) == 0 {
		return nil, nil
	}
	return posyandu[0], nil
}

func (r *Repo) GetIbuHamilByPasienID(ctx context.Context, idPasien int32) ([]*model.IbuHamil, error) {
	var hasil []*model.IbuHamil
	stmt := SELECT(IbuHamil.AllColumns).FROM(IbuHamil).WHERE(IbuHamil.IDPasien.EQ(Int32(idPasien)))
	err := pgxV5.Query(ctx, stmt, r.db, &hasil)
	if err != nil {
		return nil, err
	}
	return hasil, nil
}

func (r *Repo) GetAnakByPasienID(ctx context.Context, idPasien int32) (*model.Anak, error) {
	var hasil []*model.Anak
	stmt := SELECT(Anak.AllColumns).FROM(Anak).WHERE(Anak.IDPasien.EQ(Int32(idPasien)))
	err := pgxV5.Query(ctx, stmt, r.db, &hasil)
	if err != nil {
		return nil, err
	}
	if len(hasil) == 0 {
		return nil, nil
	}
	return hasil[0], nil
}

func (r *Repo) CreatePasien(ctx context.Context, data *model.Pasien) error {
	stmt := Pasien.INSERT(Pasien.IDPasien, Pasien.IDPosyandu, Pasien.CreatedAt, Pasien.UpdatedAt).
		VALUES(data.IDPasien, data.IDPosyandu, data.CreatedAt, data.UpdatedAt)
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) CreateIbuHamil(ctx context.Context, data *model.IbuHamil) error {
	stmt := IbuHamil.INSERT(IbuHamil.MutableColumns).
		MODEL(data)
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) CreateAnak(ctx context.Context, data *model.Anak) error {
	stmt := Anak.INSERT(Anak.IDPasien, Anak.IDIbuHamil, Anak.IDWali, Anak.NamaAnak, Anak.BeratLahir, Anak.PanjangLahir, Anak.HubunganDenganWali, Anak.CreatedAt, Anak.UpdatedAt).
		MODEL(data)
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) GetAllPaginated(ctx context.Context, page int, perPage int, q string) ([]*pasienDomain.PasienJoinRow, int, error) {
	offset := (page - 1) * perPage

	fromClause := `
		FROM pasien p
		JOIN user_account ua ON ua.id_user = p.id_pasien
		JOIN posyandu pos ON pos.id_posyandu = p.id_posyandu
		LEFT JOIN LATERAL (
			SELECT 'Ibu Hamil' AS jenis, ih.status_kehamilan::text AS sub_status
			FROM ibu_hamil ih WHERE ih.id_pasien = p.id_pasien LIMIT 1
		) ih_data ON true
		LEFT JOIN LATERAL (
			SELECT 'Anak' AS jenis, NULL AS sub_status
			FROM anak a WHERE a.id_pasien = p.id_pasien LIMIT 1
		) anak_data ON true
	`
	whereClause := ""
	args := []any{}
	argIdx := 1

	if q != "" {
		whereClause = fmt.Sprintf("WHERE (ua.nama ILIKE $%d OR ua.nik ILIKE $%d)", argIdx, argIdx)
		args = append(args, "%"+q+"%")
		argIdx++
	}

	var total int
	countQuery := "SELECT COUNT(*) " + fromClause + " " + whereClause
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	jenisCol := "COALESCE(ih_data.jenis, anak_data.jenis, 'Pasien')"
	subStatusCol := "ih_data.sub_status"

	dataQuery := fmt.Sprintf(`SELECT p.id_pasien, ua.nama, ua.nik, ua.jenis_kelamin::text, ua.tanggal_lahir::text, pos.nama_posyandu, %s AS jenis_pasien, %s AS status_kehamilan %s %s ORDER BY ua.nama ASC LIMIT $%d OFFSET $%d`,
		jenisCol, subStatusCol, fromClause, whereClause, argIdx, argIdx+1)
	args = append(args, perPage, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*pasienDomain.PasienJoinRow
	for rows.Next() {
		row := &pasienDomain.PasienJoinRow{}
		err := rows.Scan(&row.IDPasien, &row.Nama, &row.NIK, &row.JenisKelamin, &row.TanggalLahir, &row.NamaPosyandu, &row.JenisPasien, &row.StatusKehamilan)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, row)
	}

	return result, total, nil
}

func (r *Repo) Search(ctx context.Context, q string) ([]*pasienDomain.PasienJoinRow, error) {
	fromClause := `
		FROM pasien p
		JOIN user_account ua ON ua.id_user = p.id_pasien
		JOIN posyandu pos ON pos.id_posyandu = p.id_posyandu
		LEFT JOIN LATERAL (
			SELECT 'Ibu Hamil' AS jenis, ih.status_kehamilan::text AS sub_status
			FROM ibu_hamil ih WHERE ih.id_pasien = p.id_pasien LIMIT 1
		) ih_data ON true
		LEFT JOIN LATERAL (
			SELECT 'Anak' AS jenis, NULL AS sub_status
			FROM anak a WHERE a.id_pasien = p.id_pasien LIMIT 1
		) anak_data ON true
	`
	jenisCol := "COALESCE(ih_data.jenis, anak_data.jenis, 'Pasien')"
	subStatusCol := "ih_data.sub_status"

	query := fmt.Sprintf(`SELECT p.id_pasien, ua.nama, ua.nik, ua.jenis_kelamin::text, ua.tanggal_lahir::text, pos.nama_posyandu, %s AS jenis_pasien, %s AS status_kehamilan %s WHERE (ua.nama ILIKE $1 OR ua.nik ILIKE $1) ORDER BY ua.nama ASC LIMIT 20`,
		jenisCol, subStatusCol, fromClause)

	rows, err := r.db.Query(ctx, query, "%"+q+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*pasienDomain.PasienJoinRow
	for rows.Next() {
		row := &pasienDomain.PasienJoinRow{}
		err := rows.Scan(&row.IDPasien, &row.Nama, &row.NIK, &row.JenisKelamin, &row.TanggalLahir, &row.NamaPosyandu, &row.JenisPasien, &row.StatusKehamilan)
		if err != nil {
			return nil, err
		}
		result = append(result, row)
	}

	return result, nil
}

func (r *Repo) GetDetailByID(ctx context.Context, idPasien int32) (*pasienDomain.PasienDetailJoinRow, error) {
	query := `
		SELECT
			p.id_pasien,
			ua.nama,
			ua.nik,
			ua.email,
			ua.no_hp,
			ua.jenis_kelamin::text,
			ua.tanggal_lahir::text,
			ua.id_lokasi,
			pos.nama_posyandu,
			p.id_posyandu,
			COALESCE(ih_data.jenis, anak_data.jenis, 'Pasien') AS jenis_pasien,
			p.created_at::text,
			p.updated_at::text,
			ih_data.id_ibu_hamil,
			ih_data.hamil_ke,
			ih_data.bulan_mulai_hamil,
			ih_data.hpht,
			ih_data.status_kehamilan,
			anak_data.nama_anak,
			anak_data.berat_lahir,
			anak_data.panjang_lahir,
			anak_data.hubungan_dengan_wali,
			anak_data.nama_wali
		FROM pasien p
		JOIN user_account ua ON ua.id_user = p.id_pasien
		JOIN posyandu pos ON pos.id_posyandu = p.id_posyandu
		LEFT JOIN LATERAL (
			SELECT
				'Ibu Hamil' AS jenis,
				ih.id_ibu_hamil,
				ih.hamil_ke,
				ih.bulan_mulai_hamil::text AS bulan_mulai_hamil,
				ih.hpht::text AS hpht,
				ih.status_kehamilan::text AS status_kehamilan,
				NULL::text AS nama_anak,
				NULL::float8 AS berat_lahir,
				NULL::float8 AS panjang_lahir,
				NULL::text AS hubungan_dengan_wali,
				NULL::text AS nama_wali
			FROM ibu_hamil ih WHERE ih.id_pasien = p.id_pasien LIMIT 1
		) ih_data ON true
		LEFT JOIN LATERAL (
			SELECT
				'Anak' AS jenis,
				NULL::int AS id_ibu_hamil,
				NULL::int AS hamil_ke,
				NULL::text AS bulan_mulai_hamil,
				NULL::text AS hpht,
				NULL::text AS status_kehamilan,
				a.nama_anak,
				a.berat_lahir::float8,
				a.panjang_lahir::float8,
				a.hubungan_dengan_wali::text,
				w.nama AS nama_wali
			FROM anak a
			JOIN user_account w ON w.id_user = a.id_wali
			WHERE a.id_pasien = p.id_pasien LIMIT 1
		) anak_data ON true
		WHERE p.id_pasien = $1
	`

	row := &pasienDomain.PasienDetailJoinRow{}
	err := r.db.QueryRow(ctx, query, idPasien).Scan(
		&row.IDPasien, &row.Nama, &row.NIK, &row.Email, &row.NoHp,
		&row.JenisKelamin, &row.TanggalLahir, &row.IDLokasi,
		&row.NamaPosyandu, &row.IDPosyandu, &row.JenisPasien,
		&row.CreatedAt, &row.UpdatedAt,
		&row.IDIbuHamil, &row.HamilKe, &row.BulanMulaiHamil, &row.Hpht, &row.StatusKehamilan,
		&row.NamaAnak, &row.BeratLahir, &row.PanjangLahir, &row.HubunganDenganWali, &row.NamaWali,
	)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, nil
		}
		return nil, err
	}

	return row, nil
}

func (r *Repo) UpdatePasien(ctx context.Context, data *model.Pasien) error {
	data.UpdatedAt = time.Now()
	stmt := Pasien.UPDATE(Pasien.IDPosyandu, Pasien.UpdatedAt).
		MODEL(data).
		WHERE(Pasien.IDPasien.EQ(Int32(data.IDPasien)))
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) UpdateIbuHamil(ctx context.Context, data *model.IbuHamil) error {
	data.UpdatedAt = time.Now()
	stmt := IbuHamil.UPDATE(IbuHamil.HamilKe, IbuHamil.BulanMulaiHamil, IbuHamil.Hpht, IbuHamil.StatusKehamilan, IbuHamil.UpdatedAt).
		MODEL(data).
		WHERE(IbuHamil.IDIbuHamil.EQ(Int32(data.IDIbuHamil)))
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) UpdateAnak(ctx context.Context, data *model.Anak) error {
	data.UpdatedAt = time.Now()
	stmt := Anak.UPDATE(Anak.NamaAnak, Anak.BeratLahir, Anak.PanjangLahir, Anak.HubunganDenganWali, Anak.IDWali, Anak.UpdatedAt).
		MODEL(data).
		WHERE(Anak.IDPasien.EQ(Int32(data.IDPasien)))
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) DeletePasien(ctx context.Context, idPasien int32) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM anak WHERE id_pasien = $1", idPasien)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "DELETE FROM ibu_hamil WHERE id_pasien = $1", idPasien)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "DELETE FROM pasien WHERE id_pasien = $1", idPasien)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
