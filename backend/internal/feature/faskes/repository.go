package faskes

import (
	"context"

	faskesDomain "github.com/stringptr/SiGizi/backend/internal/domain/faskes"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetAll(ctx context.Context, search string) ([]*faskesDomain.FaskesItem, error) {
	query := `SELECT id_faskes, nama_faskes, tipe_faskes FROM fasilitas_kesehatan WHERE is_deleted = false`
	args := []interface{}{}
	if search != "" {
		query += ` AND nama_faskes ILIKE $1`
		args = append(args, "%"+search+"%")
	}
	query += ` ORDER BY nama_faskes ASC`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*faskesDomain.FaskesItem
	for rows.Next() {
		var item faskesDomain.FaskesItem
		if err := rows.Scan(&item.IDFaskes, &item.NamaFaskes, &item.TipeFaskes); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}

	if items == nil {
		items = []*faskesDomain.FaskesItem{}
	}
	return items, nil
}
