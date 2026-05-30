package userAccount

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
)

type Handler struct {
	Service userAccount.Service
}

func NewHandler(Service userAccount.Service) *Handler {
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

type RegisterInput struct {
	Body *userAccount.RegisterRequestDTO
}

type RegisterResponse struct {
	Body httputils.APIResponse[any]
}

func (h *Handler) Register(ctx context.Context, input *RegisterInput) (*RegisterResponse, error) {
	err := h.Service.Register(ctx, input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest("Registration Failed", err)
	}

	return &RegisterResponse{Body: httputils.Created[any](nil)}, nil
}
