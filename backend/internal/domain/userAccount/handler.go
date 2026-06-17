package userAccount

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler interface {
	GetAllUsers(ctx context.Context, input *GetAllUsersRequest) (*httputils.APIResponseOutput[*UserListData], error)
	GetUserByID(ctx context.Context, input *struct{ IDUser int32 `path:"id" minimum:"1"` }) (*httputils.APIResponseOutput[*UserDetailResponse], error)
	UpdateUser(ctx context.Context, input *UpdateUserInput) (*httputils.APIResponseOutput[*UserDetailResponse], error)
	CreateUser(ctx context.Context, input *CreateUserInput) (*httputils.APIResponseOutput[*CreateUserResponse], error)
	UpdateUserRole(ctx context.Context, input *UpdateUserRoleInput) (*httputils.APIResponseOutput[any], error)
	GetAuditLogs(ctx context.Context, input *httputils.APIRequestInput[*AuditLogFilter]) (*httputils.APIResponseOutput[*AuditLogListData], error)
}
