package repository

import (
	"database/sql"
	"fmt"
	"latihan2/app/model"
	"latihan2/database"
	"time"
)

// func GetAllAlumni() ([]model.Alumni, error) {
// 	rows, err := database.DB.Query(
// 		`SELECT id, user_id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at
// 		 FROM alumni
// 		 ORDER BY created_at ASC`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var alumniList []model.Alumni
// 	for rows.Next() {
// 		var a model.Alumni
// 		if err := rows.Scan(
// 			&a.ID, &a.UserID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
// 			&a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt,
// 		); err != nil {
// 			return nil, err
// 		}
// 		alumniList = append(alumniList, a)
// 	}
// 	return alumniList, nil
// }

func GetAlumniByID(id int, role string) (*model.Alumni, error) {

    query := `
        SELECT id, user_id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat,
               created_at, updated_at
        FROM alumni
        WHERE id = $1
    `

    var a model.Alumni
    row := database.DB.QueryRow(query, id)
    err := row.Scan(
        &a.ID, &a.UserID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
        &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt,
    )
    if err != nil {
        return nil, err
    }
    return &a, nil
}


func CreateAlumni(req model.CreateAlumniRequest) (*model.Alumni, error) {
	var newAlumni model.Alumni
	err := database.DB.QueryRow(
		`INSERT INTO alumni (user_id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		 RETURNING id, user_id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat, created_at, updated_at`,
		req.UserID, req.NIM, req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus, req.Email, req.NoTelepon, req.Alamat, time.Now(), time.Now()).Scan(
		&newAlumni.ID, &newAlumni.UserID, &newAlumni.NIM, &newAlumni.Nama, &newAlumni.Jurusan, &newAlumni.Angkatan, &newAlumni.TahunLulus,
		&newAlumni.Email, &newAlumni.NoTelepon, &newAlumni.Alamat, &newAlumni.CreatedAt, &newAlumni.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &newAlumni, nil
}

func UpdateAlumni(id int, req model.UpdateAlumniRequest) (*model.Alumni, error) {
	result, err := database.DB.Exec(
		`UPDATE alumni
		 SET nama = $1, jurusan = $2, angkatan = $3, tahun_lulus = $4, email = $5, no_telepon = $6, alamat = $7, updated_at = $8
		 WHERE id = $9`,
		req.Nama, req.Jurusan, req.Angkatan, req.TahunLulus, req.Email, req.NoTelepon, req.Alamat, time.Now(), id)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return GetAlumniByID(id, "admin")
}


func DeleteAlumni(id int) error {
	result, err := database.DB.Exec("DELETE FROM alumni WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func GetAlumniByTahunLulus(tahun int) ([]model.AlumniPekerjaanResponse, int, error) {
	db := database.DB

	rows, err := db.Query(`
		SELECT a.id, a.user_id, a.jurusan, a.tahun_lulus, p.bidang_industri, 
		       p.nama_perusahaan, p.posisi_jabatan, p.gaji_range
		FROM alumni a
		JOIN pekerjaan_alumni p ON a.id = p.alumni_id
		WHERE a.tahun_lulus = $1
		ORDER BY a.id
	`, tahun)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []model.AlumniPekerjaanResponse
	for rows.Next() {
		var res model.AlumniPekerjaanResponse
		err := rows.Scan(&res.ID, &res.Jurusan, &res.TahunLulus, &res.BidangIndustri,
			&res.NamaPerusahaan, &res.PosisiJabatan, &res.GajiRange)
		if err != nil {
			return nil, 0, err
		}
		result = append(result, res)
	}

	var total int
	err = db.QueryRow(`
		SELECT COUNT(*)
		FROM alumni a
		JOIN pekerjaan_alumni p ON a.id = p.alumni_id
		WHERE a.tahun_lulus = $1
		  AND CAST(REGEXP_REPLACE(split_part(p.gaji_range, '-', 1), '[^0-9]', '', 'g') AS INTEGER) >= 4000000
	`, tahun).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return result, total, nil
}

func GetAlumniRepo(search, sortBy, order string, limit, offset int, role string) ([]model.Alumni, error) {
    
    query := fmt.Sprintf(`
        SELECT id, user_id, nim, nama, jurusan, angkatan, tahun_lulus, email, no_telepon, alamat,
               created_at, updated_at
        FROM alumni
        WHERE (nim ILIKE $1 OR nama ILIKE $1 OR email ILIKE $1)
        ORDER BY %s %s
        LIMIT $2 OFFSET $3
    `, sortBy, order)

    rows, err := database.DB.Query(query, "%"+search+"%", limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var alumniList []model.Alumni
    for rows.Next() {
        var a model.Alumni
        if err := rows.Scan(
            &a.ID, &a.UserID, &a.NIM, &a.Nama, &a.Jurusan, &a.Angkatan, &a.TahunLulus,
            &a.Email, &a.NoTelepon, &a.Alamat, &a.CreatedAt, &a.UpdatedAt,
        ); err != nil {
            return nil, err
        }
        alumniList = append(alumniList, a)
    }
    return alumniList, nil
}

func CountAlumniRepo(search string) (int, error) {
	var total int
	countQuery := `
		SELECT COUNT(*) 
		FROM alumni 
		WHERE nim ILIKE $1 OR nama ILIKE $1 OR email ILIKE $1
	`
	err := database.DB.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}

func SoftDeleteAlumniRepo(id int) error {
    query := `UPDATE alumni SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
    res, err := database.DB.Exec(query, id)
    if err != nil {
        return err
    }

    rowsAffected, _ := res.RowsAffected()
    if rowsAffected == 0 {
        return sql.ErrNoRows
    }
    return nil
}
