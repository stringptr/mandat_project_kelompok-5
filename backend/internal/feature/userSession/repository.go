package userSession

import (
	"context"
	"errors"

	"github.com/gofrs/uuid/v5"
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

var nonDefaultMutableColumns = UserSession.MutableColumns.Except(UserSession.DefaultColumns.Except(UserSession.StatusSession))

func (r *Repo) GetByID(ctx context.Context, sessionID uuid.UUID) (*model.UserSession, error) {
	var session []*model.UserSession

	stmt := UserSession.SELECT(UserSession.AllColumns).WHERE(UserSession.IDSession.EQ(UUID(sessionID)))
	err := pgxV5.Query(ctx, stmt, r.db, session)
	if err != nil {
		return nil, err
	}

	return session[0], nil
}

func (r *Repo) Create(ctx context.Context, sessionModel *model.UserSession) error {
	stmt := UserSession.INSERT(UserSession.MutableColumns).MODEL(sessionModel)
	res, err := pgxV5.Exec(ctx, stmt, r.db)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return errors.New("fail to create session")
	}

	return nil
}
