package lokasi

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/errorutils"
)

type Service interface {
	GetLokasi(ctx context.Context, req *GetLokasiRequest) ([]*LokasiItem, *errorutils.Error)
}
