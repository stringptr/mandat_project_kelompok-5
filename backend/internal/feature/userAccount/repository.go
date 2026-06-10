package userAccount

import (
	"context"

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
