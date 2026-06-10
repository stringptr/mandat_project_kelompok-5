package pasien

import (
	"context"

	pasienDomain "github.com/stringptr/SiGizi/backend/internal/domain/pasien"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler struct {
	Service pasienDomain.Service
}

func NewHandler(service pasienDomain.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) DaftarIbuHamil(ctx context.Context, input *httputils.APIRequestInput[*pasienDomain.DaftarIbuHamilRequest]) (*httputils.APIResponseOutput[any], error) {
	err := h.Service.DaftarIbuHamil(ctx, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return &httputils.APIResponseOutput[any]{Body: httputils.Created[any](nil)}, nil
}

func (h *Handler) DaftarAnak(ctx context.Context, input *httputils.APIRequestInput[*pasienDomain.DaftarAnakRequest]) (*httputils.APIResponseOutput[any], error) {
	err := h.Service.DaftarAnak(ctx, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return &httputils.APIResponseOutput[any]{Body: httputils.Created[any](nil)}, nil
}

func (h *Handler) GetAll(ctx context.Context, input *httputils.APIRequestInput[*pasienDomain.GetAllPasienRequest]) (*httputils.APIResponseOutput[*pasienDomain.PasienListData], error) {
	res, err := h.Service.GetAll(ctx, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return &httputils.APIResponseOutput[*pasienDomain.PasienListData]{Body: httputils.OK(res)}, nil
}

func (h *Handler) Search(ctx context.Context, input *httputils.APIRequestInput[*pasienDomain.SearchPasienRequest]) (*httputils.APIResponseOutput[[]*pasienDomain.PasienListItem], error) {
	res, err := h.Service.Search(ctx, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}
	if res == nil {
		res = []*pasienDomain.PasienListItem{}
	}

	return &httputils.APIResponseOutput[[]*pasienDomain.PasienListItem]{Body: httputils.OK(res)}, nil
}

func (h *Handler) GetByID(ctx context.Context, input *struct {
	IDPasien int32 `path:"id" minimum:"1"`
},
) (*httputils.APIResponseOutput[*pasienDomain.PasienDetailResponse], error) {
	res, err := h.Service.GetByID(ctx, input.IDPasien)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return &httputils.APIResponseOutput[*pasienDomain.PasienDetailResponse]{Body: httputils.OK(res)}, nil
}

func (h *Handler) Update(ctx context.Context, input *httputils.APIRequestInput[*pasienDomain.UpdatePasienRequest]) (*httputils.APIResponseOutput[*pasienDomain.PasienDetailResponse], error) {
	res, err := h.Service.Update(ctx, input.Body)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return &httputils.APIResponseOutput[*pasienDomain.PasienDetailResponse]{Body: httputils.OK(res)}, nil
}

func (h *Handler) Delete(ctx context.Context, input *struct {
	IDPasien int32 `path:"id" minimum:"1"`
},
) (*httputils.APIResponseOutput[any], error) {
	err := h.Service.Delete(ctx, input.IDPasien)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return &httputils.APIResponseOutput[any]{Body: httputils.OK[any](nil)}, nil
}
