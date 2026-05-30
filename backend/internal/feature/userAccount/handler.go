package userAccount

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	userAccountDomain "github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Handler struct {
	Service userAccountDomain.Service
}

func NewHandler(Service userAccountDomain.Service) *Handler {
	return &Handler{Service: Service}
}

type GetAllResponse struct {
	Body httputils.APIResponse[[]*model.UserAccount]
}

func (h *Handler) GetAll(ctx context.Context, input *struct{}) (*GetAllResponse, error) {
	userAccounts, err := h.Service.GetAll(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("Internal Server Error", err)
	}

	return &GetAllResponse{Body: httputils.OK(userAccounts)}, nil
}
