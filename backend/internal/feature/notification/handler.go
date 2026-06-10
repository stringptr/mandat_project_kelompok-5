package notification

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	notificationDomain "github.com/stringptr/SiGizi/backend/internal/domain/notification"
	"github.com/stringptr/SiGizi/backend/internal/errorutils"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler struct {
	Service notificationDomain.Service
}

func NewHandler(service notificationDomain.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) GetNotifikasi(ctx context.Context, input *notificationDomain.NotifikasiListInput) (*httputils.APIResponseOutput[*notificationDomain.NotifikasiListData], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetNotifikasi(ctx, claims.IDUser, input)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return &httputils.APIResponseOutput[*notificationDomain.NotifikasiListData]{Body: httputils.OK(data)}, nil
}

func (h *Handler) GetNotifikasiDetail(ctx context.Context, input *notificationDomain.NotifikasiDetailInput) (*httputils.APIResponseOutput[*notificationDomain.NotifikasiDetail], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetNotifikasiDetail(ctx, claims.IDUser, input.IDNotifikasi)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return &httputils.APIResponseOutput[*notificationDomain.NotifikasiDetail]{Body: httputils.OK(data)}, nil
}

func (h *Handler) MarkRead(ctx context.Context, input *notificationDomain.MarkReadInput) (*httputils.APIResponseOutput[*notificationDomain.MarkReadResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.MarkRead(ctx, claims.IDUser, input.IDNotifikasi)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return &httputils.APIResponseOutput[*notificationDomain.MarkReadResponse]{Body: httputils.OK(data)}, nil
}

func (h *Handler) MarkAllRead(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*notificationDomain.MarkAllReadResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.MarkAllRead(ctx, claims.IDUser)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return &httputils.APIResponseOutput[*notificationDomain.MarkAllReadResponse]{Body: httputils.OK(data)}, nil
}

func (h *Handler) GetBidanDashboard(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*notificationDomain.BidanNotificationResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetBidanDashboard(ctx, claims.IDUser)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return &httputils.APIResponseOutput[*notificationDomain.BidanNotificationResponse]{Body: httputils.OK(data)}, nil
}

func (h *Handler) GetStatistics(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*notificationDomain.NotificationStats], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetStatistics(ctx, claims.IDUser)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return &httputils.APIResponseOutput[*notificationDomain.NotificationStats]{Body: httputils.OK(data)}, nil
}

func (h *Handler) GetActivity(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*notificationDomain.NotificationActivity], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Anda harus login untuk mengakses halaman ini.")
	}

	data, err := h.Service.GetActivity(ctx, claims.IDUser)
	if err != nil {
		return nil, errorutils.ToHumaError(err)
	}

	return &httputils.APIResponseOutput[*notificationDomain.NotificationActivity]{Body: httputils.OK(data)}, nil
}
