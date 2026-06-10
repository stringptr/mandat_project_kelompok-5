package notification

import (
	"context"
	"strings"

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

func (r *Repo) GetByID(ctx context.Context, idNotifikasi int32, idUser int32) (*model.Notifikasi, error) {
	var notif []*model.Notifikasi

	stmt := SELECT(Notifikasi.AllColumns).
		FROM(Notifikasi).
		WHERE(Notifikasi.IDNotifikasi.EQ(Int32(idNotifikasi)).
			AND(Notifikasi.IDUser.EQ(Int32(idUser))))

	err := pgxV5.Query(ctx, stmt, r.db, &notif)
	if err != nil {
		return nil, err
	}
	if len(notif) == 0 {
		return nil, nil
	}
	return notif[0], nil
}

func (r *Repo) GetByUserID(ctx context.Context, idUser int32, search string, limit int32, offset int32) ([]*model.Notifikasi, error) {
	var notif []*model.Notifikasi

	conditions := Notifikasi.IDUser.EQ(Int32(idUser))
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		conditions = conditions.AND(
			Notifikasi.Judul.ILIKE(String(searchPattern)).
				OR(Notifikasi.Pesan.ILIKE(String(searchPattern))),
		)
	}

	stmt := SELECT(Notifikasi.AllColumns).
		FROM(Notifikasi).
		WHERE(conditions).
		ORDER_BY(Notifikasi.TanggalKirim.DESC()).
		LIMIT(limit).
		OFFSET(offset)

	err := pgxV5.Query(ctx, stmt, r.db, &notif)
	if err != nil {
		return nil, err
	}
	return notif, nil
}

func (r *Repo) CountByUserID(ctx context.Context, idUser int32, search string) (int32, error) {
	var count struct {
		Count int32
	}

	conditions := Notifikasi.IDUser.EQ(Int32(idUser))
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		conditions = conditions.AND(
			Notifikasi.Judul.ILIKE(String(searchPattern)).
				OR(Notifikasi.Pesan.ILIKE(String(searchPattern))),
		)
	}

	stmt := SELECT(COUNT(Star).AS("count")).
		FROM(Notifikasi).
		WHERE(conditions)

	err := pgxV5.Query(ctx, stmt, r.db, &count)
	if err != nil {
		return 0, err
	}
	return count.Count, nil
}

func (r *Repo) CountUnreadByUserID(ctx context.Context, idUser int32) (int32, error) {
	var count struct {
		Count int32
	}

	stmt := SELECT(COUNT(Star).AS("count")).
		FROM(Notifikasi).
		WHERE(Notifikasi.IDUser.EQ(Int32(idUser)).
			AND(Notifikasi.StatusBaca.EQ(Bool(false))))

	err := pgxV5.Query(ctx, stmt, r.db, &count)
	if err != nil {
		return 0, err
	}
	return count.Count, nil
}

func (r *Repo) MarkRead(ctx context.Context, idNotifikasi int32, idUser int32) error {
	stmt := Notifikasi.UPDATE(Notifikasi.StatusBaca).
		SET(Bool(true)).
		WHERE(Notifikasi.IDNotifikasi.EQ(Int32(idNotifikasi)).
			AND(Notifikasi.IDUser.EQ(Int32(idUser))))

	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) MarkAllRead(ctx context.Context, idUser int32) (int32, error) {
	stmt := Notifikasi.UPDATE(Notifikasi.StatusBaca).
		SET(Bool(true)).
		WHERE(Notifikasi.IDUser.EQ(Int32(idUser)).
			AND(Notifikasi.StatusBaca.EQ(Bool(false))))

	res, err := pgxV5.Exec(ctx, stmt, r.db)
	if err != nil {
		return 0, err
	}
	return int32(res.RowsAffected()), nil
}

func (r *Repo) Create(ctx context.Context, data *model.Notifikasi) error {
	nonDefaultCols := Notifikasi.MutableColumns.Except(Notifikasi.DefaultColumns)

	stmt := Notifikasi.INSERT(nonDefaultCols).MODEL(data)
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	if err != nil {
		return err
	}
	return nil
}
