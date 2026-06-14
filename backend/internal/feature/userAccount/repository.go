package userAccount

import (
	"context"
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

	stmt := SELECT(UserAccount.AllColumns).FROM(UserAccount).WHERE(UserAccount.IDUser.EQ(Int32(IDUser)).AND(UserAccount.IsDeleted.EQ(Bool(false))))
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

	stmt := SELECT(UserAccount.AllColumns).FROM(UserAccount).WHERE(UserAccount.Nik.EQ(String(NIK)).AND(UserAccount.IsDeleted.EQ(Bool(false))))
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

	stmt := SELECT(UserAccount.AllColumns).FROM(UserAccount).WHERE(UserAccount.Email.EQ(String(email)).AND(UserAccount.IsDeleted.EQ(Bool(false))))
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

	stmt := SELECT(UserAccount.AllColumns).FROM(UserAccount).WHERE(UserAccount.Nik.EQ(String(NIK)).AND(UserAccount.Email.EQ(String(email)).AND(UserAccount.IsDeleted.EQ(Bool(false)))))
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
	offset := int64((page - 1) * perPage)

	fromClause := UserAccount.
		LEFT_JOIN(DinasKesehatan, DinasKesehatan.IDUser.EQ(UserAccount.IDUser)).
		LEFT_JOIN(Bidan, Bidan.IDUser.EQ(UserAccount.IDUser)).
		LEFT_JOIN(KaderPosyandu, KaderPosyandu.IDUser.EQ(UserAccount.IDUser)).
		LEFT_JOIN(Pasien, Pasien.IDPasien.EQ(UserAccount.IDUser))

	conditions := []BoolExpression{UserAccount.IsDeleted.EQ(Bool(false))}

	if q != "" {
		pattern := "%" + q + "%"
		nameILike := BoolExp(CustomExpression(UserAccount.Nama, Token("ILIKE"), String(pattern)))
		nikILike := BoolExp(CustomExpression(UserAccount.Nik, Token("ILIKE"), String(pattern)))
		conditions = append(conditions, nameILike.OR(nikILike))
	}

	if role != "" {
		switch role {
		case "Dinkes":
			conditions = append(conditions, DinasKesehatan.IDUser.IS_NOT_NULL())
		case "Bidan":
			conditions = append(conditions, Bidan.IDUser.IS_NOT_NULL())
		case "Kader":
			conditions = append(conditions, KaderPosyandu.IDUser.IS_NOT_NULL())
		case "Pasien":
			conditions = append(conditions, Pasien.IDPasien.IS_NOT_NULL())
		}
	}

	if statusVerifikasi != "" {
		conditions = append(conditions, UserAccount.StatusVerifikasi.EQ(String(statusVerifikasi)))
	}

	whereCond := conditions[0]
	for _, c := range conditions[1:] {
		whereCond = whereCond.AND(c)
	}

	var countResult struct{ Count int64 }
	countStmt := SELECT(COUNT(STAR)).FROM(fromClause).WHERE(whereCond)
	err := pgxV5.Query(ctx, countStmt, r.db, &countResult)
	if err != nil {
		return nil, 0, err
	}

	var users []*model.UserAccount
	dataStmt := SELECT(UserAccount.AllColumns).FROM(fromClause).WHERE(whereCond).
		ORDER_BY(UserAccount.CreatedAt.DESC()).
		LIMIT(int64(perPage)).
		OFFSET(offset)
	err = pgxV5.Query(ctx, dataStmt, r.db, &users)
	if err != nil {
		return nil, 0, err
	}

	return users, int(countResult.Count), nil
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
	stmt := UserAccount.UPDATE(UserAccount.StatusVerifikasi).SET(status.String()).WHERE(UserAccount.IDUser.EQ(Int32(IDUser)))

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
	stmt := UserAccount.UPDATE(UserAccount.IsDeleted, UserAccount.DeletedAt).
		SET(Bool(true), RawTimestampz("NOW()")).
		WHERE(UserAccount.IDUser.EQ(Int32(IDUser)).AND(UserAccount.IsDeleted.EQ(Bool(false))))
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}
