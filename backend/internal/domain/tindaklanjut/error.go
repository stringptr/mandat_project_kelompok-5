package tindaklanjut

import "errors"

var (
	ErrTindakLanjutNotFound   = errors.New("data tindak lanjut tidak ditemukan")
	ErrTindakLanjutNotCreated = errors.New("gagal membuat tindak lanjut")
	ErrRujukanNotFound        = errors.New("data rujukan tidak ditemukan")
	ErrStatusRujukanInvalid   = errors.New("status rujukan tidak valid")
)
