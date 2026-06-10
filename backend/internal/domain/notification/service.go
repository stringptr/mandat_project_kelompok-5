package notification

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/errorutils"
)

type Service interface {
	GetNotifikasi(ctx context.Context, idUser int32, input *NotifikasiListInput) (*NotifikasiListData, *errorutils.Error)
	GetNotifikasiDetail(ctx context.Context, idUser int32, idNotifikasi int32) (*NotifikasiDetail, *errorutils.Error)
	MarkRead(ctx context.Context, idUser int32, idNotifikasi int32) (*MarkReadResponse, *errorutils.Error)
	MarkAllRead(ctx context.Context, idUser int32) (*MarkAllReadResponse, *errorutils.Error)
	GetBidanDashboard(ctx context.Context, idUser int32) (*BidanNotificationResponse, *errorutils.Error)
	GetStatistics(ctx context.Context, idUser int32) (*NotificationStats, *errorutils.Error)
	GetActivity(ctx context.Context, idUser int32) (*NotificationActivity, *errorutils.Error)
}
