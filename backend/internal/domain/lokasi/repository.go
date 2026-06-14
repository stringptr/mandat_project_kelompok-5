package lokasi

import "context"

type Repo interface {
	GetByTipeAndParent(ctx context.Context, tipe string, bagianDari int32) ([]*LokasiItem, error)
}
