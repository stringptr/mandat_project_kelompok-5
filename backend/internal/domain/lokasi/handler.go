package lokasi

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler interface {
	GetLokasi(ctx context.Context, input *GetLokasiRequest) (*httputils.APIResponseOutput[[]*LokasiItem], error)
}
