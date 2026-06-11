package userAccount

import (
	"time"
)

type AdminGetAllResponseDTO struct {
	Email            string    `json:"email"           validate:"required,email"`
	NoHp             string    `json:"no_hp"           validate:"required,max=20"`
	Nama             string    `json:"nama"            validate:"required,min=1,max=255"`
	NIK              string    `json:"nik"             validate:"required,len=16"`
	JenisKelamin     string    `json:"jenis_kelamin"   validate:"required,oneof=Laki-Laki Perempuan"`
	TanggalLahir     time.Time `json:"tanggal_lahir"   validate:"required"`
	IDLokasi         int32     `json:"id_lokasi"       validate:"required"`
	IDPendidikan     *int32    `json:"id_pendidikan"   validate:"min=1"`
	IDPekerjaan      *int32    `json:"id_pekerjaan"    validate:"min=1"`
	IDPendapatan     *int32    `json:"id_pendapatan"   validate:"min=1"`
	JumlahTanggungan *int32    `json:"jumlah_tanggungan" validate:"min=0"`
}

type GetAllUsersRequest struct {
	Page             int     `query:"page"             minimum:"1" default:"1"`
	PerPage          int     `query:"per_page"          minimum:"1" maximum:"100" default:"20"`
	Q                string  `query:"q"                 maxLength:"255"`
	Role             string  `query:"role"              enum:"Bidan,Kader,Dinkes,Pasien"`
	StatusVerifikasi string  `query:"status_verifikasi" enum:"Pending,Aktif,Ditolak"`
}

type UserListItem struct {
	IDUser           int32    `json:"id_user"`
	Nama             string   `json:"nama"`
	NIK              string   `json:"nik"`
	Email            string   `json:"email"`
	NoHp             string   `json:"no_hp"`
	JenisKelamin     string   `json:"jenis_kelamin"`
	StatusVerifikasi string   `json:"status_verifikasi"`
	Roles            []string `json:"roles"`
	IDLokasi         int32    `json:"id_lokasi"`
	CreatedAt        string   `json:"created_at"`
	UpdatedAt        string   `json:"updated_at"`
}

type UserListData struct {
	Users     []UserListItem `json:"users"`
	TotalData int            `json:"total_data"`
	Page      int            `json:"page"`
	PerPage   int            `json:"per_page"`
}

type UserDetailResponse struct {
	IDUser           int32     `json:"id_user"`
	Email            string    `json:"email"`
	NoHp             string    `json:"no_hp"`
	Nama             string    `json:"nama"`
	NIK              string    `json:"nik"`
	JenisKelamin     string    `json:"jenis_kelamin"`
	TanggalLahir     time.Time `json:"tanggal_lahir"`
	StatusVerifikasi string    `json:"status_verifikasi"`
	IDLokasi         int32     `json:"id_lokasi"`
	IDPendidikan     *int32    `json:"id_pendidikan,omitempty"`
	IDPekerjaan      *int32    `json:"id_pekerjaan,omitempty"`
	IDPendapatan     *int32    `json:"id_pendapatan,omitempty"`
	JumlahTanggungan *int32    `json:"jumlah_tanggungan,omitempty"`
	Roles            []string  `json:"roles"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type UpdateUserInput struct {
	IDUser int32 `path:"id" minimum:"1"`
	Body   *UpdateUserRequest
}

type UpdateUserRequest struct {
	IDUser           int32      `path:"id" minimum:"1"`
	Email            *string    `json:"email,omitempty"             format:"email" maxLength:"255"`
	NoHp             *string    `json:"no_hp,omitempty"             maxLength:"20"`
	Nama             *string    `json:"nama,omitempty"              minLength:"1" maxLength:"255"`
	NIK              *string    `json:"nik,omitempty"               minLength:"16" maxLength:"16" pattern:"^[0-9]{16}$"`
	JenisKelamin     *string    `json:"jenis_kelamin,omitempty"     enum:"Laki-Laki,Perempuan"`
	TanggalLahir     *time.Time `json:"tanggal_lahir,omitempty"    format:"date-time"`
	IDLokasi         *int32     `json:"id_lokasi,omitempty"         minimum:"1"`
	IDPendidikan     *int32     `json:"id_pendidikan,omitempty"     minimum:"1"`
	IDPekerjaan      *int32     `json:"id_pekerjaan,omitempty"      minimum:"1"`
	IDPendapatan     *int32     `json:"id_pendapatan,omitempty"     minimum:"1"`
	JumlahTanggungan *int32     `json:"jumlah_tanggungan,omitempty" minimum:"0"`
}
