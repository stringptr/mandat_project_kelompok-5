package auth

import "strings"

func MapValidationError(location, message string) (string, string, bool) {
	switch {
	case strings.HasSuffix(location, "email"):
		return "ERR-VAL-01", "Masukkan alamat email yang valid", true
	case strings.HasSuffix(location, "password"):
		return "ERR-VAL-02", "Masukkan password yang lebih panjang", true
	case strings.HasSuffix(location, "nik"):
		return "ERR-VAL", "NIK yang dimasukkan tidak valid", true
	case strings.HasSuffix(location, "no_hp"):
		return "ERR-VAL-04", "Masukkan nomor telepon yang valid", true
	case strings.HasSuffix(location, "nama"):
		return "ERR-VAL-05", "Masukkan nama yang valid", true
	case strings.HasSuffix(location, "tanggal_lahir"):
		return "ERR-VAL-06", "Tanggal yang dimasukkan tidak valid", true
	case strings.HasSuffix(location, "jenis_kelamin"),
		strings.HasSuffix(location, "id_lokasi"),
		strings.HasSuffix(location, "id_pendidikan"),
		strings.HasSuffix(location, "id_pekerjaan"),
		strings.HasSuffix(location, "id_pendapatan"),
		strings.HasSuffix(location, "jumlah_tanggungan"):
		return "ERR-VAL-07", "Pilihan yang dipilih tidak valid", true
	case strings.Contains(message, "required"):
		return "ERR-VAL-03", "Lengkapi dokumen", true
	}
	return "", "", false
}
