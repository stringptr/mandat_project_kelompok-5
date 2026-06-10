package userAccount

import (
	"context"

	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler interface {
	GetAllUsers(ctx context.Context, input *httputils.APIRequestInput[*GetAllUsersRequest]) (*httputils.APIResponseOutput[*UserListData], error)
	GetUserByID(ctx context.Context, input *struct{ IDUser int32 `path:"id" minimum:"1"` }) (*httputils.APIResponseOutput[*UserDetailResponse], error)
	UpdateUser(ctx context.Context, input *httputils.APIRequestInput[*UpdateUserRequest]) (*httputils.APIResponseOutput[*UserDetailResponse], error)
}
