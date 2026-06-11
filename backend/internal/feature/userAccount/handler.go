package userAccount

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	userAccountDomain "github.com/stringptr/SiGizi/backend/internal/domain/userAccount"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type Handler struct {
	Service userAccountDomain.Service
}

func NewHandler(Service userAccountDomain.Service) *Handler {
	return &Handler{Service: Service}
}

func (h *Handler) GetAllUsers(ctx context.Context, input *httputils.APIRequestInput[*userAccountDomain.GetAllUsersRequest]) (*httputils.APIResponseOutput[*userAccountDomain.UserListData], error) {
	result, err := h.Service.GetAllUsers(ctx, input.Body)
	if err != nil {
		return nil, huma.Error500InternalServerError("Terjadi kesalahan. Silahkan dicoba kembali.", err)
	}

	return httputils.NewSuccessOutput(http.StatusOK, result, "Daftar pengguna berhasil diambil."), nil
}

func (h *Handler) GetUserByID(ctx context.Context, input *struct {
	IDUser int32 `path:"id" minimum:"1"`
},
) (*httputils.APIResponseOutput[*userAccountDomain.UserDetailResponse], error) {
	user, err := h.Service.GetUserByID(ctx, input.IDUser)
	if err != nil {
		return nil, huma.Error500InternalServerError("Terjadi kesalahan. Silahkan dicoba kembali.", err)
	}
	if user == nil {
		return nil, huma.Error404NotFound("Pengguna tidak ditemukan.")
	}

	return httputils.NewSuccessOutput(http.StatusOK, user, "Detail pengguna berhasil diambil."), nil
}

func (h *Handler) UpdateUser(ctx context.Context, input *userAccountDomain.UpdateUserInput) (*httputils.APIResponseOutput[*userAccountDomain.UserDetailResponse], error) {
	claims := httputils.GetAccessClaim(ctx)
	if claims == nil {
		return nil, huma.Error401Unauthorized("Silahkan login terlebih dahulu.")
	}

	idUser := input.IDUser
	input.Body.IDUser = idUser

	if idUser != claims.IDUser {
		isAdmin := false
		for _, r := range claims.Roles {
			if r == "ADMIN" || r == "SUPER_ADMIN" {
				isAdmin = true
				break
			}
		}
		if !isAdmin {
			return nil, huma.Error403Forbidden("Anda tidak memiliki izin untuk mengubah data pengguna lain.")
		}
	}

	err := h.Service.UpdateUser(ctx, idUser, input.Body)
	if err != nil {
		if errors.Is(err, userAccountDomain.ErrNotFound) {
			return nil, huma.Error404NotFound("Pengguna tidak ditemukan.")
		}
		return nil, huma.Error500InternalServerError("Terjadi kesalahan. Silahkan dicoba kembali.", err)
	}

	updatedUser, err := h.Service.GetUserByID(ctx, idUser)
	if err != nil {
		return nil, huma.Error500InternalServerError("Terjadi kesalahan. Silahkan dicoba kembali.", err)
	}

	return httputils.NewSuccessOutput(http.StatusOK, updatedUser, "Profil pengguna berhasil diperbarui."), nil
}
