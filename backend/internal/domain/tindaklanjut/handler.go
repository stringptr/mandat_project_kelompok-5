package tindaklanjut

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler interface {
	GetPasienTindakLanjut(ctx context.Context, input *GetPasienTindakLanjutRequest) (*httputils.APIResponseOutput[*PasienTindakLanjutData], error)
	GetDetailPasienByID(ctx context.Context, input *struct {
		IDPasien int32 `path:"id" minimum:"1"`
	}) (*httputils.APIResponseOutput[*DetailPasienTindakLanjut], error)
	CreateTindakLanjut(ctx context.Context, input *httputils.APIRequestInput[*CreateTindakLanjutRequest]) (*httputils.APIResponseOutput[*CreateTindakLanjutResponse], error)
	UpdateStatusRujukan(ctx context.Context, input *struct {
		IDRujukan int32 `path:"id" minimum:"1"`
		Body      *UpdateStatusRujukanRequest
	}) (*httputils.APIResponseOutput[*UpdateStatusRujukanResponse], error)
	GetStatusTindakLanjut(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*StatusTindakLanjutData], error)
	GetLaporanTindakLanjut(ctx context.Context, input *struct{}) (*httputils.APIResponseOutput[*LaporanTindakLanjutData], error)
	GetDetailTindakLanjutByID(ctx context.Context, input *struct {
		IDTindakLanjut int32 `path:"id" minimum:"1"`
	}) (*httputils.APIResponseOutput[*DetailTindakLanjutPasien], error)
}
