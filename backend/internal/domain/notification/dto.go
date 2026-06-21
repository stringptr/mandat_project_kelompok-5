package notification

import (
	"time"

	"github.com/stringptr/SiGizi/backend/internal/pagination"
)

type NotifikasiItem struct {
	IDNotifikasi   int32     `json:"id_notifikasi"`
	Judul          string    `json:"judul"`
	Pesan          *string   `json:"pesan,omitempty"`
	TipeNotifikasi string    `json:"tipe_notifikasi"`
	StatusBaca     bool      `json:"status_baca"`
	TanggalKirim   time.Time `json:"tanggal_kirim"`
}

type NotifikasiListData struct {
	Notifikasi []NotifikasiItem `json:"notifikasi"`
	Meta       pagination.Meta  `json:"meta"`
}

type AksiItem struct {
	Label string `json:"label"`
	URL   string `json:"url"`
}

type NotifikasiDetail struct {
	IDNotifikasi   int32     `json:"id_notifikasi"`
	Judul          string    `json:"judul"`
	Pesan          *string   `json:"pesan,omitempty"`
	TipeNotifikasi string    `json:"tipe_notifikasi"`
	StatusBaca     bool      `json:"status_baca"`
	TanggalKirim   time.Time `json:"tanggal_kirim"`
	Aksi           *AksiItem `json:"aksi,omitempty"`
}

type MarkReadResponse struct {
	IDNotifikasi int32 `json:"id_notifikasi"`
	StatusBaca   bool  `json:"status_baca"`
}

type MarkAllReadResponse struct {
	JumlahDiperbarui int32  `json:"jumlah_diperbarui"`
	Status           string `json:"status"`
}

type NotificationStats struct {
	JadwalUlang           int32 `json:"jadwal_ulang"`
	RujukanMendesak       int32 `json:"rujukan_mendesak"`
	RisikoStunting        int32 `json:"risiko_stunting"`
	NotifikasiBelumDibaca int32 `json:"notifikasi_belum_dibaca"`
}

type AktivitasItem struct {
	IDNotifikasi int32     `json:"id_notifikasi"`
	Judul        string    `json:"judul"`
	Status       string    `json:"status"`
	Timestamp    time.Time `json:"timestamp"`
}

type NotificationActivity struct {
	HariIni []AktivitasItem `json:"hari_ini"`
	Kemarin []AktivitasItem `json:"kemarin"`
}

type NotifikasiBidanStats struct {
	JadwalKontrol  int32 `json:"jadwal_kontrol"`
	RisikoStunting int32 `json:"risiko_stunting"`
	RujukanMendesak int32 `json:"rujukan_mendesak"`
}

type PasienRisikoItem struct {
	IDPasien          int32  `json:"id_pasien"`
	NamaPasien        string `json:"nama_pasien"`
	StatusGizi        string `json:"status_gizi"`
	StatusStunting    string `json:"status_stunting"`
	TanggalMonitoring string `json:"tanggal_monitoring"`
}

type JadwalMonitoringItem struct {
	IDPasien      int32  `json:"id_pasien"`
	NamaPasien    string `json:"nama_pasien"`
	JadwalKontrol string `json:"jadwal_kontrol"`
	Status        string `json:"status"`
}

type RujukanMendesakItem struct {
	IDRujukan      int32  `json:"id_rujukan"`
	NamaPasien     string `json:"nama_pasien"`
	StatusRujukan  string `json:"status_rujukan"`
	TanggalRujukan string `json:"tanggal_rujukan"`
}

type LaporanBulanan struct {
	Bulan                  string `json:"bulan"`
	JumlahPasienMonitoring int32  `json:"jumlah_pasien_monitoring"`
	JumlahPasienDirujuk    int32  `json:"jumlah_pasien_dirujuk"`
}

type BidanNotificationResponse struct {
	Statistik        NotifikasiBidanStats   `json:"statistik"`
	RisikoStunting   []PasienRisikoItem     `json:"notifikasi_risiko_stunting"`
	JadwalMonitoring []JadwalMonitoringItem `json:"jadwal_monitoring"`
	RujukanMendesak  []RujukanMendesakItem  `json:"rujukan_mendesak"`
	LaporanBulanan   LaporanBulanan         `json:"laporan_bulanan"`
}
