package auth

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	authDomain "github.com/stringptr/SiGizi/backend/internal/domain/auth"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
	"github.com/stringptr/SiGizi/backend/internal/jwtutils"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Register(ctx context.Context, input *httputils.APIRequestInput[*authDomain.RegisterRequest]) (*httputils.APIResponseOutput[any], error) {
	err := h.Service.Register(ctx, input.Body)
	if err != nil {
		return nil, huma.Error400BadRequest("Registration Failed", err)
	}

	return &httputils.APIResponseOutput[any]{Body: httputils.Created[any](nil)}, nil
}
