package lokasi

import (
	"context"

	lokasiDomain "github.com/stringptr/SiGizi/backend/internal/domain/lokasi"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetByTipeAndParent(ctx context.Context, tipe string, bagianDari int32) ([]*lokasiDomain.LokasiItem, error) {
	sql := `SELECT id_lokasi, nama_lokasi, tipe_lokasi::text, bagian_dari FROM lokasi WHERE tipe_lokasi::text = $1`

	if bagianDari > 0 {
		sql += ` AND bagian_dari = $2`
	} else {
		sql += ` AND bagian_dari IS NULL`
	}
	sql += ` ORDER BY nama_lokasi ASC`

	args := []interface{}{tipe}
	if bagianDari > 0 {
		args = append(args, bagianDari)
	}

	pgxRows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer pgxRows.Close()

	var result []*lokasiDomain.LokasiItem
	for pgxRows.Next() {
		var row lokasiDomain.LokasiItem
		err := pgxRows.Scan(&row.IDLokasi, &row.NamaLokasi, &row.TipeLokasi, &row.BagianDari)
		if err != nil {
			return nil, err
		}
		result = append(result, &row)
	}

	return result, nil
}
