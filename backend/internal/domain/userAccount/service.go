package userAccount

import (
	"context"
)

type Service interface {
	GetAllUsers(ctx context.Context, req *GetAllUsersRequest) (*UserListData, error)
	GetUserByID(ctx context.Context, idUser int32) (*UserDetailResponse, error)
	UpdateUser(ctx context.Context, idUser int32, req *UpdateUserRequest) error
}
