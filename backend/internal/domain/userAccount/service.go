package userAccount

import (
	"context"
)

type Service interface {
	GetAllUsers(ctx context.Context, req *GetAllUsersRequest) (*UserListData, error)
	GetUserByID(ctx context.Context, idUser int32) (*UserDetailResponse, error)
	UpdateUser(ctx context.Context, idUser int32, req *UpdateUserRequest) error
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
	UpdateUserRole(ctx context.Context, idUser int32, req *UpdateUserRoleRequest) error
	GetAuditLogs(ctx context.Context, filter *AuditLogFilter) (*AuditLogListData, error)
	DeleteUser(ctx context.Context, idUser int32) error
}
