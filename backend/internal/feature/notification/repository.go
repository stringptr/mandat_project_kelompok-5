package notification

import (
	"context"
	"strings"
	"time"

	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
	. "github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/table"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/go-jet/jet/v2/pgxV5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetByID(ctx context.Context, idNotifikasi int32, idUser int32) (*model.Notifikasi, error) {
	var notif []*model.Notifikasi

	stmt := SELECT(Notifikasi.AllColumns).
		FROM(Notifikasi).
		WHERE(Notifikasi.IDNotifikasi.EQ(Int32(idNotifikasi)).
			AND(Notifikasi.IDUser.EQ(Int32(idUser))))

	err := pgxV5.Query(ctx, stmt, r.db, &notif)
	if err != nil {
		return nil, err
	}
	if len(notif) == 0 {
		return nil, nil
	}
	return notif[0], nil
}

func (r *Repo) GetByUserID(ctx context.Context, idUser int32, search string, limit int32, offset int32) ([]*model.Notifikasi, error) {
	var notif []*model.Notifikasi

	conditions := Notifikasi.IDUser.EQ(Int32(idUser))
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		conditions = conditions.AND(
			Notifikasi.Judul.ILIKE(String(searchPattern)).
				OR(Notifikasi.Pesan.ILIKE(String(searchPattern))),
		)
	}

	stmt := SELECT(Notifikasi.AllColumns).
		FROM(Notifikasi).
		WHERE(conditions).
		ORDER_BY(Notifikasi.TanggalKirim.DESC()).
		LIMIT(limit).
		OFFSET(offset)

	err := pgxV5.Query(ctx, stmt, r.db, &notif)
	if err != nil {
		return nil, err
	}
	return notif, nil
}

func (r *Repo) CountByUserID(ctx context.Context, idUser int32, search string) (int32, error) {
	var count struct {
		Count int32
	}

	conditions := Notifikasi.IDUser.EQ(Int32(idUser))
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		conditions = conditions.AND(
			Notifikasi.Judul.ILIKE(String(searchPattern)).
				OR(Notifikasi.Pesan.ILIKE(String(searchPattern))),
		)
	}

	stmt := SELECT(COUNT(Star).AS("count")).
		FROM(Notifikasi).
		WHERE(conditions)

	err := pgxV5.Query(ctx, stmt, r.db, &count)
	if err != nil {
		return 0, err
	}
	return count.Count, nil
}

func (r *Repo) CountUnreadByUserID(ctx context.Context, idUser int32) (int32, error) {
	var count struct {
		Count int32
	}

	stmt := SELECT(COUNT(Star).AS("count")).
		FROM(Notifikasi).
		WHERE(Notifikasi.IDUser.EQ(Int32(idUser)).
			AND(Notifikasi.StatusBaca.EQ(Bool(false))))

	err := pgxV5.Query(ctx, stmt, r.db, &count)
	if err != nil {
		return 0, err
	}
	return count.Count, nil
}

func (r *Repo) MarkRead(ctx context.Context, idNotifikasi int32, idUser int32) error {
	stmt := Notifikasi.UPDATE(Notifikasi.StatusBaca).
		SET(Bool(true)).
		WHERE(Notifikasi.IDNotifikasi.EQ(Int32(idNotifikasi)).
			AND(Notifikasi.IDUser.EQ(Int32(idUser))))

	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) MarkAllRead(ctx context.Context, idUser int32) (int32, error) {
	stmt := Notifikasi.UPDATE(Notifikasi.StatusBaca).
		SET(Bool(true)).
		WHERE(Notifikasi.IDUser.EQ(Int32(idUser)).
			AND(Notifikasi.StatusBaca.EQ(Bool(false))))

	res, err := pgxV5.Exec(ctx, stmt, r.db)
	if err != nil {
		return 0, err
	}
	return int32(res.RowsAffected()), nil
}

func (r *Repo) Create(ctx context.Context, data *model.Notifikasi) error {
	nonDefaultCols := Notifikasi.MutableColumns.Except(Notifikasi.DefaultColumns)

	stmt := Notifikasi.INSERT(nonDefaultCols).MODEL(data)
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) getBidanIDByUserID(ctx context.Context, idUser int32) (int32, error) {
	var bidans []*model.Bidan

	stmt := SELECT(Bidan.IDBidan).
		FROM(Bidan).
		WHERE(Bidan.IDUser.EQ(Int32(idUser))).
		LIMIT(1)

	err := pgxV5.Query(ctx, stmt, r.db, &bidans)
	if err != nil {
		return 0, err
	}
	if len(bidans) == 0 {
		return 0, nil
	}
	return bidans[0].IDBidan, nil
}

// ---------------------------------------------------------------------------
// Private helpers — query result rows
// ---------------------------------------------------------------------------

type riskStuntingRow struct {
	IDPasien          int32     `sql:"alias:id_pasien"`
	NamaPasien        string    `sql:"alias:nama_pasien"`
	StatusGizi        string    `sql:"alias:status_gizi"`
	StatusStunting    string    `sql:"alias:status_stunting"`
	TanggalMonitoring time.Time `sql:"alias:tanggal_monitoring"`
}

type jadwalMonitoringRow struct {
	IDPasien      int32     `sql:"alias:id_pasien"`
	NamaPasien    string    `sql:"alias:nama_pasien"`
	JadwalKontrol time.Time `sql:"alias:jadwal_kontrol"`
	Status        string    `sql:"alias:status"`
}

type rujukanMendesakRow struct {
	IDRujukan      int32     `sql:"alias:id_rujukan"`
	NamaPasien     string    `sql:"alias:nama_pasien"`
	StatusRujukan  string    `sql:"alias:status_rujukan"`
	TanggalRujukan time.Time `sql:"alias:tanggal_rujukan"`
}

type laporanBulananRow struct {
	Bulan                  string `sql:"alias:bulan"`
	JumlahPasienMonitoring int32  `sql:"alias:jumlah_pasien_monitoring"`
	JumlahPasienDirujuk    int32  `sql:"alias:jumlah_pasien_dirujuk"`
}

// ---------------------------------------------------------------------------
// Aktivitas
// ---------------------------------------------------------------------------

func (r *Repo) getTodayNotifikasi(ctx context.Context, idUser int32) ([]*model.Notifikasi, error) {
	var notif []*model.Notifikasi

	stmt := SELECT(Notifikasi.AllColumns).
		FROM(Notifikasi).
		WHERE(Notifikasi.IDUser.EQ(Int32(idUser)).
			AND(Notifikasi.TanggalKirim.GT_EQ(Raw("CURRENT_DATE")))).
		ORDER_BY(Notifikasi.TanggalKirim.DESC())

	err := pgxV5.Query(ctx, stmt, r.db, &notif)
	if err != nil {
		return nil, err
	}
	return notif, nil
}

func (r *Repo) getYesterdayNotifikasi(ctx context.Context, idUser int32) ([]*model.Notifikasi, error) {
	var notif []*model.Notifikasi

	stmt := SELECT(Notifikasi.AllColumns).
		FROM(Notifikasi).
		WHERE(Notifikasi.IDUser.EQ(Int32(idUser)).
			AND(Notifikasi.TanggalKirim.GT_EQ(Raw("CURRENT_DATE - INTERVAL '1 day'"))).
			AND(Notifikasi.TanggalKirim.LT(Raw("CURRENT_DATE")))).
		ORDER_BY(Notifikasi.TanggalKirim.DESC())

	err := pgxV5.Query(ctx, stmt, r.db, &notif)
	if err != nil {
		return nil, err
	}
	return notif, nil
}

// ---------------------------------------------------------------------------
// Statistik — counts per bidan
// ---------------------------------------------------------------------------

func (r *Repo) countBidanJadwalKontrol(ctx context.Context, idBidan int32) (int32, error) {
	var count struct {
		Count int32
	}

	stmt := SELECT(COUNT(Star).AS("count")).
		FROM(TindakLanjut).
		WHERE(TindakLanjut.IDBidan.EQ(Int32(idBidan)).
			AND(TindakLanjut.JadwalKontrol.GT_EQ(Raw("CURRENT_DATE"))))

	err := pgxV5.Query(ctx, stmt, r.db, &count)
	if err != nil {
		return 0, err
	}
	return count.Count, nil
}

func (r *Repo) countBidanRisikoStunting(ctx context.Context, idBidan int32) (int32, error) {
	var count struct {
		Count int32
	}

	stmt := SELECT(COUNT(DISTINCT(JadwalImunisasi.IDPasien)).AS("count")).
		FROM(
			HasilPemeriksaan.
			INNER_JOIN(JadwalImunisasi, HasilPemeriksaan.IDJadwalImunisasi.EQ(JadwalImunisasi.IDImunisasi)).
			INNER_JOIN(Pasien, JadwalImunisasi.IDPasien.EQ(Pasien.IDPasien)).
			INNER_JOIN(Posyandu, Pasien.IDPosyandu.EQ(Posyandu.IDPosyandu)),
		).
		WHERE(Posyandu.IDBidan.EQ(Int32(idBidan)).
			AND(HasilPemeriksaan.StatusStunting.IN(
				String("Berisiko Stunting"),
				String("Stunting"),
				String("Stunting Berat"),
			)))

	err := pgxV5.Query(ctx, stmt, r.db, &count)
	if err != nil {
		return 0, err
	}
	return count.Count, nil
}

func (r *Repo) countBidanRujukanMendesak(ctx context.Context, idBidan int32) (int32, error) {
	var count struct {
		Count int32
	}

	stmt := SELECT(COUNT(Star).AS("count")).
		FROM(
			Rujukan.
			INNER_JOIN(TindakLanjut, Rujukan.IDTindakLanjut.EQ(TindakLanjut.IDTindakLanjut)),
		).
		WHERE(TindakLanjut.IDBidan.EQ(Int32(idBidan)).
			AND(Rujukan.StatusRujukan.IN(
				String("Diajukan"),
				String("Diproses"),
			)))

	err := pgxV5.Query(ctx, stmt, r.db, &count)
	if err != nil {
		return 0, err
	}
	return count.Count, nil
}

// ---------------------------------------------------------------------------
// Bidan dashboard — list data
// ---------------------------------------------------------------------------

func (r *Repo) getBidanRisikoStuntingList(ctx context.Context, idBidan int32) ([]*riskStuntingRow, error) {
	var rows []*riskStuntingRow

	stmt := SELECT(
		Pasien.IDPasien.AS("id_pasien"),
		Anak.NamaAnak.AS("nama_pasien"),
		HasilPemeriksaan.StatusGizi.AS("status_gizi"),
		HasilPemeriksaan.StatusStunting.AS("status_stunting"),
		HasilPemeriksaan.CreatedAt.AS("tanggal_monitoring"),
	).
		FROM(
			HasilPemeriksaan.
			INNER_JOIN(JadwalImunisasi, HasilPemeriksaan.IDJadwalImunisasi.EQ(JadwalImunisasi.IDImunisasi)).
			INNER_JOIN(Pasien, JadwalImunisasi.IDPasien.EQ(Pasien.IDPasien)).
			INNER_JOIN(Anak, Pasien.IDPasien.EQ(Anak.IDPasien)).
			INNER_JOIN(Posyandu, Pasien.IDPosyandu.EQ(Posyandu.IDPosyandu)),
		).
		WHERE(Posyandu.IDBidan.EQ(Int32(idBidan)).
			AND(HasilPemeriksaan.StatusStunting.IN(
				String("Berisiko Stunting"),
				String("Stunting"),
				String("Stunting Berat"),
			))).
		ORDER_BY(HasilPemeriksaan.CreatedAt.DESC())

	err := pgxV5.Query(ctx, stmt, r.db, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repo) getBidanJadwalMonitoringList(ctx context.Context, idBidan int32) ([]*jadwalMonitoringRow, error) {
	var rows []*jadwalMonitoringRow

	stmt := SELECT(
		Pasien.IDPasien.AS("id_pasien"),
		Anak.NamaAnak.AS("nama_pasien"),
		TindakLanjut.JadwalKontrol.AS("jadwal_kontrol"),
		TindakLanjut.StatusPasien.AS("status"),
	).
		FROM(
			TindakLanjut.
			INNER_JOIN(HasilPemeriksaan, TindakLanjut.IDHasilPemeriksaan.EQ(HasilPemeriksaan.IDHasilPemeriksaan)).
			INNER_JOIN(JadwalImunisasi, HasilPemeriksaan.IDJadwalImunisasi.EQ(JadwalImunisasi.IDImunisasi)).
			INNER_JOIN(Pasien, JadwalImunisasi.IDPasien.EQ(Pasien.IDPasien)).
			INNER_JOIN(Anak, Pasien.IDPasien.EQ(Anak.IDPasien)),
		).
		WHERE(TindakLanjut.IDBidan.EQ(Int32(idBidan)).
			AND(TindakLanjut.JadwalKontrol.GT_EQ(Raw("CURRENT_DATE")))).
		ORDER_BY(TindakLanjut.JadwalKontrol.ASC())

	err := pgxV5.Query(ctx, stmt, r.db, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repo) getBidanRujukanMendesakList(ctx context.Context, idBidan int32) ([]*rujukanMendesakRow, error) {
	var rows []*rujukanMendesakRow

	stmt := SELECT(
		Rujukan.IDRujukan.AS("id_rujukan"),
		Anak.NamaAnak.AS("nama_pasien"),
		Rujukan.StatusRujukan.AS("status_rujukan"),
		Rujukan.TanggalRujukan.AS("tanggal_rujukan"),
	).
		FROM(
			Rujukan.
			INNER_JOIN(TindakLanjut, Rujukan.IDTindakLanjut.EQ(TindakLanjut.IDTindakLanjut)).
			INNER_JOIN(HasilPemeriksaan, TindakLanjut.IDHasilPemeriksaan.EQ(HasilPemeriksaan.IDHasilPemeriksaan)).
			INNER_JOIN(JadwalImunisasi, HasilPemeriksaan.IDJadwalImunisasi.EQ(JadwalImunisasi.IDImunisasi)).
			INNER_JOIN(Pasien, JadwalImunisasi.IDPasien.EQ(Pasien.IDPasien)).
			INNER_JOIN(Anak, Pasien.IDPasien.EQ(Anak.IDPasien)),
		).
		WHERE(TindakLanjut.IDBidan.EQ(Int32(idBidan)).
			AND(Rujukan.StatusRujukan.IN(
				String("Diajukan"),
				String("Diproses"),
			))).
		ORDER_BY(Rujukan.TanggalRujukan.DESC())

	err := pgxV5.Query(ctx, stmt, r.db, &rows)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *Repo) getBidanLaporanBulanan(ctx context.Context, idBidan int32) (*laporanBulananRow, error) {
	var row laporanBulananRow

	stmt := SELECT(
		Raw("TO_CHAR(CURRENT_DATE, 'FMMonth YYYY')").AS("bulan"),
		COUNT(DISTINCT(JadwalImunisasi.IDPasien)).AS("jumlah_pasien_monitoring"),
		COUNT(DISTINCT(Rujukan.IDRujukan)).AS("jumlah_pasien_dirujuk"),
	).
		FROM(
			TindakLanjut.
			LEFT_JOIN(HasilPemeriksaan, TindakLanjut.IDHasilPemeriksaan.EQ(HasilPemeriksaan.IDHasilPemeriksaan)).
			LEFT_JOIN(JadwalImunisasi, HasilPemeriksaan.IDJadwalImunisasi.EQ(JadwalImunisasi.IDImunisasi)).
			LEFT_JOIN(Rujukan, Rujukan.IDTindakLanjut.EQ(TindakLanjut.IDTindakLanjut)),
		).
		WHERE(TindakLanjut.IDBidan.EQ(Int32(idBidan)).
			AND(TindakLanjut.CreatedAt.GT_EQ(Raw("DATE_TRUNC('month', CURRENT_DATE)"))).
			AND(TindakLanjut.CreatedAt.LT(Raw("DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month'"))))

	err := pgxV5.Query(ctx, stmt, r.db, &row)
	if err != nil {
		return nil, err
	}
	return &row, nil
}
