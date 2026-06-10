package notification

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler interface {
	GetNotifikasi(ctx context.Context, input *NotifikasiListInput) (*httputils.APIResponseOutput[*NotifikasiListData], error)
	GetNotifikasiDetail(ctx context.Context, input *NotifikasiDetailInput) (*httputils.APIResponseOutput[*NotifikasiDetail], error)
	MarkRead(ctx context.Context, input *MarkReadInput) (*httputils.APIResponseOutput[*MarkReadResponse], error)
	MarkAllRead(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*MarkAllReadResponse], error)
	GetBidanDashboard(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*BidanNotificationResponse], error)
	GetStatistics(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*NotificationStats], error)
	GetActivity(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*NotificationActivity], error)
}

type NotifikasiListInput struct {
	Page    int32  `query:"page"    default:"1"`
	PerPage int32  `query:"per_page" default:"15" maximum:"100"`
	Search  string `query:"q,omitempty"`
}

type NotifikasiDetailInput struct {
	IDNotifikasi int32 `path:"id" minimum:"1"`
}

type MarkReadInput struct {
	IDNotifikasi int32 `path:"id" minimum:"1"`
}
