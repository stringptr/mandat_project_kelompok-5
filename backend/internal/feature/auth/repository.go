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
	stmt := SELECT(
		Raw("EXISTS(SELECT 1 FROM dinas_kesehatan WHERE id_user = #1 AND is_deleted = false)", RawArgs{"#1": idUser}).AS("is_dinkes"),
		Raw("EXISTS(SELECT 1 FROM bidan WHERE id_user = #1 AND is_deleted = false)", RawArgs{"#1": idUser}).AS("is_bidan"),
		Raw("EXISTS(SELECT 1 FROM kader_posyandu WHERE id_user = #1 AND is_deleted = false)", RawArgs{"#1": idUser}).AS("is_kader"),
		Raw("EXISTS(SELECT 1 FROM pasien WHERE id_pasien = #1 AND is_deleted = false)", RawArgs{"#1": idUser}).AS("is_pasien"),
		Raw("EXISTS(SELECT 1 FROM ibu_hamil WHERE id_pasien = #1 AND is_deleted = false)", RawArgs{"#1": idUser}).AS("is_ibu_hamil"),
		Raw("EXISTS(SELECT 1 FROM anak WHERE id_pasien = #1 AND is_deleted = false)", RawArgs{"#1": idUser}).AS("is_anak"),
	)

	var result struct {
		IsDinkes   bool
		IsBidan    bool
		IsKader    bool
		IsPasien   bool
		IsIbuHamil bool
		IsAnak     bool
	}
	err := pgxV5.Query(ctx, stmt, r.db, &result)
	if err != nil {
		return nil, err
	}

	isDinkes, isBidan, isKader, isPasien, isIbuHamil, isAnak := result.IsDinkes, result.IsBidan, result.IsKader, result.IsPasien, result.IsIbuHamil, result.IsAnak

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

	cond := UserAccount.IsDeleted.EQ(Bool(false))
	if email != "" {
		cond = cond.AND(UserAccount.Email.EQ(String(email)))
	}
	if nik != "" {
		cond = cond.AND(UserAccount.Nik.EQ(String(nik)))
	}

	stmt := SELECT(UserAccount.AllColumns).
		FROM(UserAccount).
		WHERE(cond).
		LIMIT(2)

	err := pgxV5.Query(ctx, stmt, r.db, &users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 || len(users) > 1 {
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

func (r *Repo) CreateRoleRecord(ctx context.Context, idUser int32, role string, noStr string, wilayahKerja int32, noSk string, idPosyandu int32) error {
	now := time.Now()
	switch role {
	case "Bidan":
		stmt := Bidan.INSERT(Bidan.IDUser, Bidan.NoStr, Bidan.WilayahKerja, Bidan.CreatedAt, Bidan.UpdatedAt).MODEL(&model.Bidan{
			IDUser:       idUser,
			NoStr:        noStr,
			WilayahKerja: wilayahKerja,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
		_, err := pgxV5.Exec(ctx, stmt, r.db)
		return err
	case "Kader":
		stmt := KaderPosyandu.INSERT(KaderPosyandu.IDUser, KaderPosyandu.NoSk, KaderPosyandu.IDPosyandu, KaderPosyandu.CreatedAt, KaderPosyandu.UpdatedAt).MODEL(&model.KaderPosyandu{
			IDUser:     idUser,
			NoSk:       noSk,
			IDPosyandu: idPosyandu,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
		_, err := pgxV5.Exec(ctx, stmt, r.db)
		return err
	case "Dinkes":
		stmt := DinasKesehatan.INSERT(DinasKesehatan.IDUser, DinasKesehatan.CreatedAt, DinasKesehatan.UpdatedAt).MODEL(&model.DinasKesehatan{
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

func (r *Repo) GetDinkesUserIDs(ctx context.Context) ([]int32, error) {
	var rows []*struct {
		IDUser int32
	}
	sql := `
		SELECT dk.id_user
		FROM dinas_kesehatan dk
		JOIN user_account ua ON ua.id_user = dk.id_user
		WHERE dk.is_deleted = false AND ua.is_deleted = false
	`

	err := pgxV5.Query(ctx, RawStatement(sql), r.db, &rows)
	if err != nil {
		return nil, err
	}
	ids := make([]int32, len(rows))
	for i, row := range rows {
		ids[i] = row.IDUser
	}
	return ids, nil
}

func (r *Repo) DeleteRoleRecords(ctx context.Context, idUser int32) error {
	bidanStmt := Bidan.UPDATE(Bidan.IsDeleted, Bidan.DeletedAt).
		SET(Bool(true), RawTimestampz("NOW()")).
		WHERE(Bidan.IDUser.EQ(Int32(idUser)))
	_, err := pgxV5.Exec(ctx, bidanStmt, r.db)
	if err != nil {
		return err
	}

	kaderStmt := KaderPosyandu.UPDATE(KaderPosyandu.IsDeleted, KaderPosyandu.DeletedAt).
		SET(Bool(true), RawTimestampz("NOW()")).
		WHERE(KaderPosyandu.IDUser.EQ(Int32(idUser)))
	_, err = pgxV5.Exec(ctx, kaderStmt, r.db)
	if err != nil {
		return err
	}

	dinkesStmt := DinasKesehatan.UPDATE(DinasKesehatan.IsDeleted, DinasKesehatan.DeletedAt).
		SET(Bool(true), RawTimestampz("NOW()")).
		WHERE(DinasKesehatan.IDUser.EQ(Int32(idUser)))
	_, err = pgxV5.Exec(ctx, dinkesStmt, r.db)
	return err
}
