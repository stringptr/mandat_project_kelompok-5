package auth

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetRoles(ctx context.Context, idUser int32) ([]string, error) {
	query := `
        SELECT
            EXISTS(SELECT 1 FROM dinas_kesehatan WHERE id_user = $1) as is_dinkes,
            EXISTS(SELECT 1 FROM bidan WHERE id_user = $1) as is_bidan,
            EXISTS(SELECT 1 FROM kader_posyandu WHERE id_user = $1) as is_kader,
            EXISTS(SELECT 1 FROM pasien WHERE id_pasien = $1) as is_pasien,
            EXISTS(SELECT 1 FROM ibu_hamil WHERE id_pasien = $1) as is_ibu_hamil,
            EXISTS(SELECT 1 FROM anak WHERE id_pasien = $1) as is_anak
    `
	var isDinkes, isBidan, isKader, isPasien, isIbuHamil, isAnak bool
	err := r.db.QueryRow(ctx, query, idUser).Scan(&isDinkes, &isBidan, &isKader, &isPasien, &isIbuHamil, &isAnak)
	if err != nil {
		return nil, err
	}

	roles := []string{"USER"}
	if isDinkes {
		roles = append(roles, "ADMIN_DINKES")
	}
	if isBidan {
		roles = append(roles, "BIDAN")
	}
	if isKader {
		roles = append(roles, "KADER")
	}
	if isKader || isBidan {
		roles = append(roles, "ADMIN")
	}
	if isPasien {
		roles = append(roles, "PASIEN")
	}
	if isIbuHamil {
		roles = append(roles, "IBU_HAMIL")
	}
	if isAnak {
		roles = append(roles, "ANAK")
	}

	return roles, nil
}
