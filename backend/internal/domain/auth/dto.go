package auth

import (
	"net/http"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/stringptr/SiGizi/backend/internal/httputils"
)

type RegisterRequest struct {
	Email            string    `json:"email"           format:"email"`
	Password         string    `json:"password"        minLength:"8" maxLength:"255"`
	NoHp             string    `json:"no_hp"           maxLength:"20"`
	Nama             string    `json:"nama"            minLength:"1" maxLength:"255"`
	NIK              string    `json:"nik"             minLength:"16" maxLength:"16" pattern:"^[0-9]{16}$"`
	JenisKelamin     string    `json:"jenis_kelamin"   enum:"Laki-Laki,Perempuan"`
	TanggalLahir     time.Time `json:"tanggal_lahir"   format:"date-time"`
	IDLokasi         int32     `json:"id_lokasi"       minimum:"1"`
	IDPendidikan     *int32    `json:"id_pendidikan,omitempty"     minimum:"1"`
	IDPekerjaan      *int32    `json:"id_pekerjaan,omitempty"      minimum:"1"`
	IDPendapatan     *int32    `json:"id_pendapatan,omitempty"     minimum:"1"`
	JumlahTanggungan *int32    `json:"jumlah_tanggungan,omitempty" minimum:"0"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type LoginRequest struct {
	Email    string `json:"email"           format:"email"`
	Password string `json:"password"        minLength:"8" maxLength:"255"`
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

type VerifyUserRequest struct {
	IDUser int32  `path:"id_user" minimum:"1"`
	Status string `json:"status" enum:"Aktif,Ditolak"`
}
