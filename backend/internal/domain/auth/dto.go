package auth

import (
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type RegisterRequest struct {
	Email            string    `json:"email"           validate:"required,email"`
	Password         string    `json:"password"        validate:"required,max=255"`
	NoHp             string    `json:"no_hp"           validate:"required,max=20"`
	Nama             string    `json:"nama"            validate:"required,min=1,max=255"`
	NIK              string    `json:"nik"             validate:"required,len=16"`
	JenisKelamin     string    `json:"jenis_kelamin"   validate:"required,oneof=Laki-Laki Perempuan"`
	TanggalLahir     time.Time `json:"tanggal_lahir"   validate:"required"`
	IDLokasi         int32     `json:"id_lokasi"       validate:"required"`
	IDPendidikan     *int32    `json:"id_pendidikan,omitempty"   validate:"min=1"`
	IDPekerjaan      *int32    `json:"id_pekerjaan,omitempty"    validate:"min=1"`
	IDPendapatan     *int32    `json:"id_pendapatan,omitempty"   validate:"min=1"`
	JumlahTanggungan *int32    `json:"jumlah_tanggungan,omitempty" validate:"min=0"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken           string    `json:"access_token"`
	RefreshToken          uuid.UUID `json:"refresh_token"`
	AccessTokenExpiresIn  int64     `json:"access_token_expires_in"`
	RefreshTokenExpiresIn int64     `json:"refresh_token_expires_in"`
}

type AuthOutput struct {
	Body      httputils.APIResponse[*AuthResponse]
	SetCookie []http.Cookie `header:"Set-Cookie"`
}

type LogoutOutput struct {
	Body      httputils.APIResponse[any]
	SetCookie []http.Cookie `header:"Set-Cookie"`
}
