package pemeriksaan

import "errors"

var (
	ErrPemeriksaanNotFound          = errors.New("data pemeriksaan tidak ditemukan")
	ErrJadwalImunisasiNotFound      = errors.New("data jadwal imunisasi tidak ditemukan")
	ErrPemeriksaanNotCreated        = errors.New("gagal membuat data pemeriksaan")
	ErrPemeriksaanNotUpdated        = errors.New("gagal memperbarui data pemeriksaan")
	ErrPemeriksaanNotDeleted        = errors.New("gagal menghapus data pemeriksaan")
	ErrPemeriksaanNotVerified       = errors.New("gagal memverifikasi data pemeriksaan")
	ErrPemeriksaanNotPending        = errors.New("tidak ada data pemeriksaan yang menunggu verifikasi")
)
