package pasien

import (
	"context"
	"strings"
	"time"

	pasienDomain "github.com/stringptr/SiGizi/backend/internal/domain/pasien"
	"github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/model"
	. "github.com/stringptr/SiGizi/backend/internal/infrastructure/jet/imunisasi/public/table"

	"github.com/go-jet/jet/v2/pgxV5"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetUserByID(ctx context.Context, idUser int32) (*model.UserAccount, error) {
	var users []*model.UserAccount
	stmt := SELECT(UserAccount.AllColumns).FROM(UserAccount).WHERE(UserAccount.IDUser.EQ(Int32(idUser)).AND(UserAccount.IsDeleted.EQ(Bool(false))))
	err := pgxV5.Query(ctx, stmt, r.db, &users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

func (r *Repo) GetPasienByUserID(ctx context.Context, idUser int32) (*model.Pasien, error) {
	var pasien []*model.Pasien
	stmt := SELECT(Pasien.AllColumns).FROM(Pasien).WHERE(Pasien.IDPasien.EQ(Int32(idUser)).AND(Pasien.IsDeleted.EQ(Bool(false))))
	err := pgxV5.Query(ctx, stmt, r.db, &pasien)
	if err != nil {
		return nil, err
	}
	if len(pasien) == 0 {
		return nil, nil
	}
	return pasien[0], nil
}

func (r *Repo) GetPosyanduByID(ctx context.Context, idPosyandu int32) (*model.Posyandu, error) {
	var posyandu []*model.Posyandu
	stmt := SELECT(Posyandu.AllColumns).FROM(Posyandu).WHERE(Posyandu.IDPosyandu.EQ(Int32(idPosyandu)).AND(Posyandu.IsDeleted.EQ(Bool(false))))
	err := pgxV5.Query(ctx, stmt, r.db, &posyandu)
	if err != nil {
		return nil, err
	}
	if len(posyandu) == 0 {
		return nil, nil
	}
	return posyandu[0], nil
}

func (r *Repo) GetIbuHamilByPasienID(ctx context.Context, idPasien int32) ([]*model.IbuHamil, error) {
	var hasil []*model.IbuHamil
	stmt := SELECT(IbuHamil.AllColumns).FROM(IbuHamil).WHERE(IbuHamil.IDPasien.EQ(Int32(idPasien)).AND(IbuHamil.IsDeleted.EQ(Bool(false))))
	err := pgxV5.Query(ctx, stmt, r.db, &hasil)
	if err != nil {
		return nil, err
	}
	return hasil, nil
}

func (r *Repo) GetAnakByPasienID(ctx context.Context, idPasien int32) (*model.Anak, error) {
	var hasil []*model.Anak
	stmt := SELECT(Anak.AllColumns).FROM(Anak).WHERE(Anak.IDPasien.EQ(Int32(idPasien)).AND(Anak.IsDeleted.EQ(Bool(false))))
	err := pgxV5.Query(ctx, stmt, r.db, &hasil)
	if err != nil {
		return nil, err
	}
	if len(hasil) == 0 {
		return nil, nil
	}
	return hasil[0], nil
}

func (r *Repo) CreatePasien(ctx context.Context, data *model.Pasien) error {
	stmt := Pasien.INSERT(Pasien.IDPasien, Pasien.IDPosyandu, Pasien.CreatedAt, Pasien.UpdatedAt).
		VALUES(data.IDPasien, data.IDPosyandu, data.CreatedAt, data.UpdatedAt)
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) CreateIbuHamil(ctx context.Context, data *model.IbuHamil) error {
	stmt := IbuHamil.INSERT(IbuHamil.MutableColumns).
		MODEL(data)
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) CreateAnak(ctx context.Context, data *model.Anak) error {
	stmt := Anak.INSERT(Anak.IDPasien, Anak.IDIbuHamil, Anak.IDWali, Anak.NamaAnak, Anak.BeratLahir, Anak.PanjangLahir, Anak.HubunganDenganWali, Anak.CreatedAt, Anak.UpdatedAt).
		MODEL(data)
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) CheckPasienOwnership(ctx context.Context, idPasien int32, idUser int32) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM pasien p
			LEFT JOIN anak a ON a.id_pasien = p.id_pasien AND a.is_deleted = false
			WHERE p.id_pasien = #1
			  AND p.is_deleted = false
			  AND (p.id_pasien = #2 OR a.id_wali = #2)
		) AS owned
	`
	var result struct{ Owned bool }
	err := pgxV5.Query(ctx, RawStatement(query, RawArgs{"#1": idPasien, "#2": idUser}), r.db, &result)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return false, nil
		}
		return false, err
	}
	return result.Owned, nil
}

func (r *Repo) GetAllPaginated(ctx context.Context, page int, perPage int, q string) ([]*pasienDomain.PasienJoinRow, int, error) {
	offset := int64((page - 1) * perPage)

	p := Pasien.AS("p")
	ua := UserAccount.AS("ua")
	pos := Posyandu.AS("pos")

	fromClause := p.
		INNER_JOIN(ua, ua.IDUser.EQ(p.IDPasien).AND(ua.IsDeleted.EQ(Bool(false)))).
		INNER_JOIN(pos, pos.IDPosyandu.EQ(p.IDPosyandu).AND(pos.IsDeleted.EQ(Bool(false))))

	conditions := []BoolExpression{p.IsDeleted.EQ(Bool(false))}

	if q != "" {
		pattern := "%" + q + "%"
		nameILike := BoolExp(CustomExpression(ua.Nama, Token("ILIKE"), String(pattern)))
		nikILike := BoolExp(CustomExpression(ua.Nik, Token("ILIKE"), String(pattern)))
		conditions = append(conditions, nameILike.OR(nikILike))
	}

	whereCond := conditions[0]
	for _, c := range conditions[1:] {
		whereCond = whereCond.AND(c)
	}

	var countResult struct{ Count int64 }
	countStmt := SELECT(COUNT(STAR)).FROM(fromClause).WHERE(whereCond)
	err := pgxV5.Query(ctx, countStmt, r.db, &countResult)
	if err != nil {
		return nil, 0, err
	}

	jenisIbuHamil := Raw("(SELECT 'Ibu Hamil' FROM ibu_hamil ih WHERE ih.id_pasien = p.id_pasien AND ih.is_deleted = false LIMIT 1)")
	jenisAnak := Raw("(SELECT 'Anak' FROM anak a WHERE a.id_pasien = p.id_pasien AND a.is_deleted = false LIMIT 1)")
	subStatus := Raw("(SELECT ih.status_kehamilan::text FROM ibu_hamil ih WHERE ih.id_pasien = p.id_pasien AND ih.is_deleted = false LIMIT 1)")

	var result []*pasienDomain.PasienJoinRow
	dataStmt := SELECT(
		p.IDPasien,
		ua.Nama,
		ua.Nik,
		Raw("ua.jenis_kelamin::text"),
		Raw("ua.tanggal_lahir::text"),
		pos.NamaPosyandu,
		COALESCE(jenisIbuHamil, jenisAnak, String("Pasien")).AS("jenis_pasien"),
		subStatus.AS("status_kehamilan"),
	).FROM(fromClause).WHERE(whereCond).
		ORDER_BY(ua.Nama.ASC()).
		LIMIT(int64(perPage)).
		OFFSET(offset)
	err = pgxV5.Query(ctx, dataStmt, r.db, &result)
	if err != nil {
		return nil, 0, err
	}

	return result, int(countResult.Count), nil
}

func (r *Repo) GetAllPaginatedByUser(ctx context.Context, page int, perPage int, q string, idUser int32) ([]*pasienDomain.PasienJoinRow, int, error) {
	offset := int64((page - 1) * perPage)

	p := Pasien.AS("p")
	ua := UserAccount.AS("ua")
	pos := Posyandu.AS("pos")
	anakAlias := Anak.AS("a")

	fromClause := p.
		INNER_JOIN(ua, ua.IDUser.EQ(p.IDPasien).AND(ua.IsDeleted.EQ(Bool(false)))).
		INNER_JOIN(pos, pos.IDPosyandu.EQ(p.IDPosyandu).AND(pos.IsDeleted.EQ(Bool(false)))).
		LEFT_JOIN(anakAlias, anakAlias.IDPasien.EQ(p.IDPasien).AND(anakAlias.IsDeleted.EQ(Bool(false))))

	ownershipCond := p.IDPasien.EQ(Int32(idUser)).OR(anakAlias.IDWali.EQ(Int32(idUser)))
	conditions := []BoolExpression{p.IsDeleted.EQ(Bool(false)), ownershipCond}

	if q != "" {
		pattern := "%" + q + "%"
		nameILike := BoolExp(CustomExpression(ua.Nama, Token("ILIKE"), String(pattern)))
		nikILike := BoolExp(CustomExpression(ua.Nik, Token("ILIKE"), String(pattern)))
		conditions = append(conditions, nameILike.OR(nikILike))
	}

	whereCond := conditions[0]
	for _, c := range conditions[1:] {
		whereCond = whereCond.AND(c)
	}

	var countResult struct{ Count int64 }
	countStmt := SELECT(COUNT(STAR)).FROM(fromClause).WHERE(whereCond)
	err := pgxV5.Query(ctx, countStmt, r.db, &countResult)
	if err != nil {
		return nil, 0, err
	}

	jenisIbuHamil := Raw("(SELECT 'Ibu Hamil' FROM ibu_hamil ih WHERE ih.id_pasien = p.id_pasien AND ih.is_deleted = false LIMIT 1)")
	jenisAnak := Raw("(SELECT 'Anak' FROM anak a WHERE a.id_pasien = p.id_pasien AND a.is_deleted = false LIMIT 1)")
	subStatus := Raw("(SELECT ih.status_kehamilan::text FROM ibu_hamil ih WHERE ih.id_pasien = p.id_pasien AND ih.is_deleted = false LIMIT 1)")

	var result []*pasienDomain.PasienJoinRow
	dataStmt := SELECT(
		p.IDPasien,
		ua.Nama,
		ua.Nik,
		Raw("ua.jenis_kelamin::text"),
		Raw("ua.tanggal_lahir::text"),
		pos.NamaPosyandu,
		COALESCE(jenisIbuHamil, jenisAnak, String("Pasien")).AS("jenis_pasien"),
		subStatus.AS("status_kehamilan"),
	).FROM(fromClause).WHERE(whereCond).
		ORDER_BY(ua.Nama.ASC()).
		LIMIT(int64(perPage)).
		OFFSET(offset)
	err = pgxV5.Query(ctx, dataStmt, r.db, &result)
	if err != nil {
		return nil, 0, err
	}

	return result, int(countResult.Count), nil
}

func (r *Repo) Search(ctx context.Context, q string) ([]*pasienDomain.PasienJoinRow, error) {
	p := Pasien.AS("p")
	ua := UserAccount.AS("ua")
	pos := Posyandu.AS("pos")

	fromClause := p.
		INNER_JOIN(ua, ua.IDUser.EQ(p.IDPasien).AND(ua.IsDeleted.EQ(Bool(false)))).
		INNER_JOIN(pos, pos.IDPosyandu.EQ(p.IDPosyandu).AND(pos.IsDeleted.EQ(Bool(false))))

	pattern := "%" + q + "%"
	whereCond := p.IsDeleted.EQ(Bool(false)).AND(
		BoolExp(CustomExpression(ua.Nama, Token("ILIKE"), String(pattern))).
			OR(BoolExp(CustomExpression(ua.Nik, Token("ILIKE"), String(pattern)))),
	)

	jenisIbuHamil := Raw("(SELECT 'Ibu Hamil' FROM ibu_hamil ih WHERE ih.id_pasien = p.id_pasien AND ih.is_deleted = false LIMIT 1)")
	jenisAnak := Raw("(SELECT 'Anak' FROM anak a WHERE a.id_pasien = p.id_pasien AND a.is_deleted = false LIMIT 1)")
	subStatus := Raw("(SELECT ih.status_kehamilan::text FROM ibu_hamil ih WHERE ih.id_pasien = p.id_pasien AND ih.is_deleted = false LIMIT 1)")

	var result []*pasienDomain.PasienJoinRow
	err := pgxV5.Query(ctx,
		SELECT(
			p.IDPasien,
			ua.Nama,
			ua.Nik,
			Raw("ua.jenis_kelamin::text"),
			Raw("ua.tanggal_lahir::text"),
			pos.NamaPosyandu,
			COALESCE(jenisIbuHamil, jenisAnak, String("Pasien")).AS("jenis_pasien"),
			subStatus.AS("status_kehamilan"),
		).FROM(fromClause).WHERE(whereCond).ORDER_BY(ua.Nama.ASC()).LIMIT(20),
		r.db, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repo) GetDetailByID(ctx context.Context, idPasien int32) (*pasienDomain.PasienDetailJoinRow, error) {
	dataSQL := `
		SELECT
			p.id_pasien,
			ua.nama,
			ua.nik,
			ua.email,
			ua.no_hp,
			ua.jenis_kelamin::text,
			ua.tanggal_lahir::text,
			ua.id_lokasi,
			pos.nama_posyandu,
			p.id_posyandu,
			COALESCE(ih_data.jenis, anak_data.jenis, 'Pasien') AS jenis_pasien,
			p.created_at::text,
			p.updated_at::text,
			ih_data.id_ibu_hamil,
			ih_data.hamil_ke,
			ih_data.bulan_mulai_hamil,
			ih_data.hpht,
			ih_data.status_kehamilan,
			anak_data.nama_anak,
			anak_data.berat_lahir,
			anak_data.panjang_lahir,
			anak_data.hubungan_dengan_wali,
			anak_data.nama_wali
		FROM pasien p
		JOIN user_account ua ON ua.id_user = p.id_pasien
		JOIN posyandu pos ON pos.id_posyandu = p.id_posyandu
		LEFT JOIN LATERAL (
			SELECT
				'Ibu Hamil' AS jenis,
				ih.id_ibu_hamil,
				ih.hamil_ke,
				ih.bulan_mulai_hamil::text AS bulan_mulai_hamil,
				ih.hpht::text AS hpht,
				ih.status_kehamilan::text AS status_kehamilan,
				NULL::text AS nama_anak,
				NULL::float8 AS berat_lahir,
				NULL::float8 AS panjang_lahir,
				NULL::text AS hubungan_dengan_wali,
				NULL::text AS nama_wali
			FROM ibu_hamil ih WHERE ih.id_pasien = p.id_pasien LIMIT 1
		) ih_data ON true
		LEFT JOIN LATERAL (
			SELECT
				'Anak' AS jenis,
				NULL::int AS id_ibu_hamil,
				NULL::int AS hamil_ke,
				NULL::text AS bulan_mulai_hamil,
				NULL::text AS hpht,
				NULL::text AS status_kehamilan,
				a.nama_anak,
				a.berat_lahir::float8,
				a.panjang_lahir::float8,
				a.hubungan_dengan_wali::text,
				w.nama AS nama_wali
			FROM anak a
			JOIN user_account w ON w.id_user = a.id_wali
			WHERE a.id_pasien = p.id_pasien LIMIT 1
		) anak_data ON true
		WHERE p.id_pasien = #1 AND p.is_deleted = false
	`

	var row pasienDomain.PasienDetailJoinRow
	err := pgxV5.Query(ctx, RawStatement(dataSQL, RawArgs{"#1": idPasien}), r.db, &row)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

func (r *Repo) UpdatePasien(ctx context.Context, data *model.Pasien) error {
	data.UpdatedAt = time.Now()
	stmt := Pasien.UPDATE(Pasien.IDPosyandu, Pasien.UpdatedAt).
		MODEL(data).
		WHERE(Pasien.IDPasien.EQ(Int32(data.IDPasien)))
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) UpdateIbuHamil(ctx context.Context, data *model.IbuHamil) error {
	data.UpdatedAt = time.Now()
	stmt := IbuHamil.UPDATE(IbuHamil.HamilKe, IbuHamil.BulanMulaiHamil, IbuHamil.Hpht, IbuHamil.StatusKehamilan, IbuHamil.UpdatedAt).
		MODEL(data).
		WHERE(IbuHamil.IDIbuHamil.EQ(Int32(data.IDIbuHamil)))
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) UpdateAnak(ctx context.Context, data *model.Anak) error {
	data.UpdatedAt = time.Now()
	stmt := Anak.UPDATE(Anak.NamaAnak, Anak.BeratLahir, Anak.PanjangLahir, Anak.HubunganDenganWali, Anak.IDWali, Anak.UpdatedAt).
		MODEL(data).
		WHERE(Anak.IDPasien.EQ(Int32(data.IDPasien)))
	_, err := pgxV5.Exec(ctx, stmt, r.db)
	return err
}

func (r *Repo) DeletePasien(ctx context.Context, idPasien int32) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	anakStmt := Anak.UPDATE(Anak.IsDeleted, Anak.DeletedAt).
		SET(Bool(true), RawTimestampz("NOW()")).
		WHERE(Anak.IDPasien.EQ(Int32(idPasien)))
	_, err = pgxV5.Exec(ctx, anakStmt, tx)
	if err != nil {
		return err
	}

	ibuHamilStmt := IbuHamil.UPDATE(IbuHamil.IsDeleted, IbuHamil.DeletedAt).
		SET(Bool(true), RawTimestampz("NOW()")).
		WHERE(IbuHamil.IDPasien.EQ(Int32(idPasien)))
	_, err = pgxV5.Exec(ctx, ibuHamilStmt, tx)
	if err != nil {
		return err
	}

	pasienStmt := Pasien.UPDATE(Pasien.IsDeleted, Pasien.DeletedAt).
		SET(Bool(true), RawTimestampz("NOW()")).
		WHERE(Pasien.IDPasien.EQ(Int32(idPasien)))
	_, err = pgxV5.Exec(ctx, pasienStmt, tx)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
