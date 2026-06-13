package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
	. "github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/table"

	"github.com/go-jet/jet/v2/pgxV5"
	. "github.com/go-jet/jet/v2/postgres"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetRoles(ctx context.Context, idUser int32) ([]string, error) {
	query := `
        SELECT
            EXISTS(SELECT 1 FROM dinas_kesehatan WHERE id_user = $1) as is_dinkes,
            EXISTS(SELECT 1 FROM bidan WHERE id_user = $1) as is_bidan,
            EXISTS(SELECT 1 FROM kader_posyandu WHERE id_user = $1) as is_kader,
            EXISTS(SELECT 1 FROM pasien WHERE id_pasien = $1) as is_pasien,
            EXISTS(SELECT 1 FROM ibu_hamil WHERE id_pasien = $1) as is_ibu_hamil,
            EXISTS(SELECT 1 FROM anak WHERE id_pasien = $1) as is_anak
    `
	var isDinkes, isBidan, isKader, isPasien, isIbuHamil, isAnak bool
	err := r.db.QueryRow(ctx, query, idUser).Scan(&isDinkes, &isBidan, &isKader, &isPasien, &isIbuHamil, &isAnak)
	if err != nil {
		return nil, err
	}

	roles := []string{"USER"}
	if isDinkes {
		roles = append(roles, "DINKES")
		roles = append(roles, "SUPER_ADMIN")
		roles = append(roles, "ADMIN")
	}
	if isBidan {
		roles = append(roles, "BIDAN")
	}
	if isKader {
		roles = append(roles, "KADER")
	}
	if isKader || isBidan {
		roles = append(roles, "ADMIN")
	}
	if isPasien {
		roles = append(roles, "PASIEN")
	}
	if isIbuHamil {
		roles = append(roles, "IBU_HAMIL")
	}
	if isAnak {
		roles = append(roles, "ANAK")
	}

	return roles, nil
}

func (r *Repo) GetByEmailNIK(ctx context.Context, email string, nik string) (*model.UserAccount, error) {
	var users []*model.UserAccount

	stmt := SELECT(UserAccount.AllColumns).
		FROM(UserAccount).
		WHERE(UserAccount.Email.EQ(String(email)).
			AND(UserAccount.Nik.EQ(String(nik)))).
		LIMIT(1)

	err := pgxV5.Query(ctx, stmt, r.db, &users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

func (r *Repo) CreateUser(ctx context.Context, data *model.UserAccount) (int32, error) {
	nonDefaultCols := UserAccount.MutableColumns.Except(UserAccount.DefaultColumns)
	var users []*model.UserAccount

	stmt := UserAccount.INSERT(nonDefaultCols).MODEL(data).RETURNING(UserAccount.AllColumns)
	err := pgxV5.Query(ctx, stmt, r.db, &users)
	if err != nil {
		return 0, err
	}
	if len(users) == 0 {
		return 0, errors.New("gagal membuat akun")
	}
	return users[0].IDUser, nil
}

func (r *Repo) CreateRoleRecord(ctx context.Context, idUser int32, role string, noStr string, wilayahKerja int32, noSk string) error {
	now := time.Now()
	switch role {
	case "Bidan":
		stmt := Bidan.INSERT(Bidan.MutableColumns).MODEL(&model.Bidan{
			IDUser:       idUser,
			NoStr:        noStr,
			WilayahKerja: wilayahKerja,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		_, err := pgxV5.Exec(ctx, stmt, r.db)
		return err
	case "Kader":
		stmt := KaderPosyandu.INSERT(KaderPosyandu.MutableColumns).MODEL(&model.KaderPosyandu{
			IDUser:     idUser,
			NoSk:       noSk,
			IDPosyandu: 0,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
		_, err := pgxV5.Exec(ctx, stmt, r.db)
		return err
	case "Dinkes":
		stmt := DinasKesehatan.INSERT(DinasKesehatan.MutableColumns).MODEL(&model.DinasKesehatan{
			IDUser:    idUser,
			CreatedAt: now,
			UpdatedAt: now,
		})
		_, err := pgxV5.Exec(ctx, stmt, r.db)
		return err
	default:
		return nil
	}
}

func (r *Repo) DeleteRoleRecords(ctx context.Context, idUser int32) error {
	_, err := r.db.Exec(ctx, "DELETE FROM bidan WHERE id_user = $1", idUser)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, "DELETE FROM kader_posyandu WHERE id_user = $1", idUser)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, "DELETE FROM dinas_kesehatan WHERE id_user = $1", idUser)
	return err
}
