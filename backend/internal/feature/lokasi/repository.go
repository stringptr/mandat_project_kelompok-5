package lokasi

import (
	"context"

	lokasiDomain "github.com/stringptr/SiGizi/backend/internal/domain/lokasi"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/enum"
	. "github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/table"

	"github.com/go-jet/jet/v2/pgxV5"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

var tipeLokasiMap = map[string]StringExpression{
	"Provinsi":  enum.TipeLokasi.Provinsi,
	"Kabupaten": enum.TipeLokasi.Kabupaten,
	"Kota":      enum.TipeLokasi.Kota,
	"Kecamatan": enum.TipeLokasi.Kecamatan,
	"Kelurahan": enum.TipeLokasi.Kelurahan,
}

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetByTipeAndParent(ctx context.Context, tipe string, bagianDari int32) ([]*lokasiDomain.LokasiItem, error) {
	conditions := []BoolExpression{
		Lokasi.TipeLokasi.EQ(tipeLokasiMap[tipe]),
	}

	if bagianDari > 0 {
		conditions = append(conditions, Lokasi.BagianDari.EQ(Int32(bagianDari)))
	} else {
		conditions = append(conditions, Lokasi.BagianDari.IS_NULL())
	}

	whereCond := conditions[0]
	for _, c := range conditions[1:] {
		whereCond = whereCond.AND(c)
	}

	var result []*lokasiDomain.LokasiItem
	stmt := SELECT(
		Lokasi.IDLokasi,
		Lokasi.NamaLokasi,
		Lokasi.TipeLokasi,
		Lokasi.BagianDari,
	).FROM(Lokasi).WHERE(whereCond).ORDER_BY(Lokasi.NamaLokasi.ASC())

	err := pgxV5.Query(ctx, stmt, r.db, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
