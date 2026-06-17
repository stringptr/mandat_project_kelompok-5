package tindaklanjut

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/errorutils"
)

type Service interface {
	GetPasienTindakLanjut(ctx context.Context, req *GetPasienTindakLanjutRequest) (*PasienTindakLanjutData, *errorutils.Error)
	GetDetailPasienByID(ctx context.Context, idPasien int32) (*DetailPasienTindakLanjut, *errorutils.Error)
	CreateTindakLanjut(ctx context.Context, idBidan int32, req *CreateTindakLanjutRequest) (*CreateTindakLanjutResponse, *errorutils.Error)
	UpdateStatusRujukan(ctx context.Context, idRujukan int32, req *UpdateStatusRujukanRequest) (*UpdateStatusRujukanResponse, *errorutils.Error)
	GetStatusTindakLanjut(ctx context.Context) (*StatusTindakLanjutData, *errorutils.Error)
	GetLaporanTindakLanjut(ctx context.Context) (*LaporanTindakLanjutData, *errorutils.Error)
	GetDetailTindakLanjutByID(ctx context.Context, idTindakLanjut int32) (*DetailTindakLanjutPasien, *errorutils.Error)
}
