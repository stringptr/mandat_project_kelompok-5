package userAccount

import (
	"context"
	"fmt"
	"strings"
	"time"

	userAccountDomain "github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
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

var MutableNonDefaultColumns ColumnList = UserAccount.MutableColumns.Except(UserAccount.DefaultColumns)

func (r *Repo) GetByID(ctx context.Context, IDUser int32) (*model.UserAccount, error) {
	var user []*model.UserAccount

	stmt := SELECT(UserAccount.AllColumns).FROM(UserAccount).WHERE(UserAccount.IDUser.EQ(Int32(IDUser)))
	err := pgxV5.Query(ctx, stmt, r.db, &user)
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, nil
	}

	return user[0], nil
}

func (r *Repo) GetByNIK(ctx context.Context, NIK string) (*model.UserAccount, error) {
	var user []*model.UserAccount

	stmt := SELECT(UserAccount.AllColumns).FROM(UserAccount).WHERE(UserAccount.Nik.EQ(String(NIK)))
	err := pgxV5.Query(ctx, stmt, r.db, &user)
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, nil
	}

	return user[0], nil
}

func (r *Repo) GetByEmail(ctx context.Context, email string) (*model.UserAccount, error) {
	var user []*model.UserAccount

	stmt := SELECT(UserAccount.AllColumns).FROM(UserAccount).WHERE(UserAccount.Email.EQ(String(email)))
	err := pgxV5.Query(ctx, stmt, r.db, &user)
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, nil
	}

	return user[0], nil
}

func (r *Repo) GetByNIKEmail(ctx context.Context, NIK string, email string) (*model.UserAccount, error) {
	var user []*model.UserAccount

	stmt := SELECT(UserAccount.AllColumns).FROM(UserAccount).WHERE(UserAccount.Nik.EQ(String(NIK)).AND(UserAccount.Email.EQ(String(email))))
	err := pgxV5.Query(ctx, stmt, r.db, &user)
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, nil
	}

	return user[0], nil
}

func (r *Repo) GetAll(ctx context.Context) ([]*model.UserAccount, error) {
	var users []*model.UserAccount

	stmt := SELECT(UserAccount.AllColumns).FROM(UserAccount)
	err := pgxV5.Query(ctx, stmt, r.db, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repo) GetAllPaginated(ctx context.Context, page int, perPage int, q string, role string, statusVerifikasi string) ([]*model.UserAccount, int, error) {
	offset := (page - 1) * perPage

	baseQuery := `
		FROM user_account ua
		LEFT JOIN dinas_kesehatan dk ON dk.id_user = ua.id_user
		LEFT JOIN bidan b ON b.id_user = ua.id_user
		LEFT JOIN kader_posyandu kp ON kp.id_user = ua.id_user
		LEFT JOIN pasien p ON p.id_pasien = ua.id_user
	`
	var conditions []string
	var args []any
	argIdx := 1

	if q != "" {
		conditions = append(conditions, fmt.Sprintf("(ua.nama ILIKE $%d OR ua.nik ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+q+"%")
		argIdx++
	}

	if role != "" {
		switch role {
		case "Dinkes":
			conditions = append(conditions, "dk.id_user IS NOT NULL")
		case "Bidan":
			conditions = append(conditions, "b.id_user IS NOT NULL")
		case "Kader":
			conditions = append(conditions, "kp.id_user IS NOT NULL")
		case "Pasien":
			conditions = append(conditions, "p.id_pasien IS NOT NULL")
		}
	}

	if statusVerifikasi != "" {
		conditions = append(conditions, fmt.Sprintf("ua.status_verifikasi = $%d", argIdx))
		args = append(args, statusVerifikasi)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	countQuery := "SELECT COUNT(*) " + baseQuery + " " + whereClause
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	dataQuery := fmt.Sprintf(`SELECT ua.id_user, ua.email, ua.password, ua.no_hp, ua.status_verifikasi, ua.nama, ua.nik, ua.jenis_kelamin, ua.tanggal_lahir, ua.id_lokasi, ua.id_pendidikan, ua.id_pekerjaan, ua.id_pendapatan, ua.jumlah_tanggungan, ua.created_at, ua.updated_at %s %s ORDER BY ua.created_at DESC LIMIT $%d OFFSET $%d`,
		baseQuery, whereClause, argIdx, argIdx+1)
	args = append(args, perPage, offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*model.UserAccount
	for rows.Next() {
		var u model.UserAccount
		err := rows.Scan(
			&u.IDUser, &u.Email, &u.Password, &u.NoHp,
			&u.StatusVerifikasi, &u.Nama, &u.Nik,
			&u.JenisKelamin, &u.TanggalLahir, &u.IDLokasi,
			&u.IDPendidikan, &u.IDPekerjaan, &u.IDPendapatan,
			&u.JumlahTanggungan, &u.CreatedAt, &u.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &u)
	}

	return users, total, nil
}

func (r *Repo) Create(ctx context.Context, userAccModel *model.UserAccount) error {
	stmt := UserAccount.INSERT(MutableNonDefaultColumns).MODEL(userAccModel)
	res, err := pgxV5.Exec(ctx, stmt, r.db)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return userAccountDomain.ErrNotCreated
	}
	return nil
}

func (r *Repo) Update(ctx context.Context, dataModel *model.UserAccount) error {
	dataModel.UpdatedAt = time.Now()

	stmt := UserAccount.UPDATE(
		UserAccount.Email,
		UserAccount.NoHp,
		UserAccount.Nama,
		UserAccount.Nik,
		UserAccount.JenisKelamin,
		UserAccount.TanggalLahir,
		UserAccount.IDLokasi,
		UserAccount.IDPendidikan,
		UserAccount.IDPekerjaan,
		UserAccount.IDPendapatan,
		UserAccount.JumlahTanggungan,
		UserAccount.UpdatedAt,
	).MODEL(dataModel).WHERE(UserAccount.IDUser.EQ(Int32(dataModel.IDUser)))

	res, err := pgxV5.Exec(ctx, stmt, r.db)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return userAccountDomain.ErrNotUpdated
	}

	return nil
}

func (r *Repo) UpdateStatusVerifikasi(ctx context.Context, IDUser int32, status model.StatusVerifikasi) error {
	stmt := UserAccount.UPDATE(UserAccount.StatusVerifikasi, UserAccount.UpdatedAt).
		SET(String(status.String()), Raw("CURRENT_TIMESTAMP")).
		WHERE(UserAccount.IDUser.EQ(Int32(IDUser)))

	res, err := pgxV5.Exec(ctx, stmt, r.db)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return userAccountDomain.ErrNotUpdated
	}

	return nil
}

func (r *Repo) DeleteByID(ctx context.Context, IDUser int32) error {
	stmt := UserAccount.DELETE().WHERE(UserAccount.IDUser.EQ(Int32(IDUser)))
	res, err := pgxV5.Exec(ctx, stmt, r.db)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return userAccountDomain.ErrNotDeleted
	}

	return nil
}
