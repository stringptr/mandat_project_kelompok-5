package userAccount

import (
	"context"
	"errors"
	"fmt"

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

func (r *Repo) GetByID(ctx context.Context, IDUser int32) (*model.UserAccount, error) {
	var user []*model.UserAccount

	stmt := SELECT(UserAccount.AllColumns).FROM(UserAccount).WHERE(UserAccount.IDUser.EQ(Int32(IDUser)))
	err := pgxV5.Query(ctx, stmt, r.db, &user)
	if err != nil {
		return nil, err
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

func (r *Repo) GetAll(ctx context.Context) ([]*model.UserAccount, error) {
	var users []*model.UserAccount

	stmt := SELECT(UserAccount.AllColumns).FROM(UserAccount)
	err := pgxV5.Query(ctx, stmt, r.db, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repo) Create(ctx context.Context, userAccModel *model.UserAccount) error {
	stmt := UserAccount.INSERT(UserAccount.AllColumns).MODEL(userAccModel)
	res, err := pgxV5.Exec(ctx, stmt, r.db)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("user account creation failed")
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
		return fmt.Errorf("user account dengan id %d cannot be found", IDUser)
	}

	return nil
}

func (r *Repo) UpdateStatusVerifikasiByID(ctx context.Context, IDUser int32, statusVerifikasi bool) error {
	stmt := UserAccount.UPDATE(UserAccount.StatusVerifikasi).SET(statusVerifikasi).WHERE(UserAccount.IDUser.EQ(Int32(IDUser)))
	res, err := pgxV5.Exec(ctx, stmt, r.db)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("user account with iduser of %d cannot be found", IDUser)
	}

	return nil
}
