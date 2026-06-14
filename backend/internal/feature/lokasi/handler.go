package lokasi

import (
	"context"

	lokasiDomain "github.com/stringptr/SiGizi/backend/internal/domain/lokasi"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler struct {
	Service lokasiDomain.Service
}

func NewHandler(service lokasiDomain.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) GetLokasi(ctx context.Context, input *lokasiDomain.GetLokasiRequest) (*httputils.APIResponseOutput[[]*lokasiDomain.LokasiItem], error) {
	res, err := h.Service.GetLokasi(ctx, input)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return httputils.NewOKOutput(res), nil
}
