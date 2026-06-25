package faskes

import (
	"context"

	faskesDomain "github.com/stringptr/SiGizi/backend/internal/domain/faskes"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler struct {
	Service faskesDomain.Service
}

func NewHandler(service faskesDomain.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) GetFaskes(ctx context.Context, input *faskesDomain.GetFaskesRequest) (*httputils.APIResponseOutput[[]*faskesDomain.FaskesItem], error) {
	res, err := h.Service.GetFaskes(ctx, input)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	return httputils.NewOKOutput(res), nil
}
