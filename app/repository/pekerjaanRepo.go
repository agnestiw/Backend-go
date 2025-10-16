package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"latihan2/app/model"
	"latihan2/database"
	"log"
	"time"
)

var (
	ErrPekerjaanNotFound = errors.New("pekerjaan not found")
	ErrForbidden         = errors.New("forbidden access")
)

func GetPekerjaanByIDRepo(id int) (*model.Pekerjaan, error) {
	var p model.Pekerjaan
	query := `
        SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
               lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
               status_pekerjaan, deskripsi_pekerjaan, is_delete, delete_by, deleted_at, created_at
        FROM pekerjaan_alumni
        WHERE id = $1
    `
	err := database.DB.QueryRow(query, id).Scan(
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan,
		&p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja,
		&p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.Deskripsi,
		&p.IsDelete, &p.DeletedBy, &p.DeletedAt, &p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetPekerjaanByAlumniID(alumniID int) ([]model.Pekerjaan, error) {
	rows, err := database.DB.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
			   lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
			   status_pekerjaan, deskripsi_pekerjaan, is_delete, delete_by, deleted_at, created_at
		FROM pekerjaan_alumni
		WHERE alumni_id = $1
		ORDER BY tanggal_mulai_kerja DESC`, alumniID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		if err := rows.Scan(
			&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan,
			&p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja,
			&p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.Deskripsi,
			&p.IsDelete, &p.DeletedBy, &p.DeletedAt, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func CreatePekerjaan(req model.CreatePekerjaanRequest) (*model.Pekerjaan, error) {
	var newPekerjaan model.Pekerjaan
	query := `
	INSERT INTO pekerjaan_alumni 
	(alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
	 lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
	 status_pekerjaan, deskripsi_pekerjaan, is_delete, created_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,FALSE,$11)
	RETURNING id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
	          lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
	          status_pekerjaan, deskripsi_pekerjaan, is_delete, created_at
`

	var tglSelesaiVal sql.NullTime
	if req.TanggalSelesaiKerja != "" {
		if t, err := time.Parse("2006-01-02", req.TanggalSelesaiKerja); err == nil {
			tglSelesaiVal = sql.NullTime{Time: t, Valid: true}
		}
	}

	err := database.DB.QueryRow(query,
		req.AlumniID, req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri,
		req.LokasiKerja, req.GajiRange, req.TanggalMulaiKerja, tglSelesaiVal,
		req.StatusPekerjaan, req.DeskripsiPekerjaan, time.Now(),
	).Scan(
		&newPekerjaan.ID, &newPekerjaan.AlumniID,
		&newPekerjaan.NamaPerusahaan, &newPekerjaan.PosisiJabatan, &newPekerjaan.BidangIndustri,
		&newPekerjaan.LokasiKerja, &newPekerjaan.GajiRange, &newPekerjaan.TanggalMulaiKerja,
		&newPekerjaan.TanggalSelesaiKerja, &newPekerjaan.StatusPekerjaan,
		&newPekerjaan.Deskripsi, &newPekerjaan.IsDelete, &newPekerjaan.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &newPekerjaan, nil
}

func UpdatePekerjaan(id int, req model.UpdatePekerjaanRequest) (*model.Pekerjaan, error) {
	var tglSelesaiVal sql.NullTime
	if req.TanggalSelesaiKerja != "" {
		if t, err := time.Parse("2006-01-02", req.TanggalSelesaiKerja); err == nil {
			tglSelesaiVal = sql.NullTime{Time: t, Valid: true}
		}
	}

	query := `
		UPDATE pekerjaan_alumni
		SET nama_perusahaan=$1, posisi_jabatan=$2, bidang_industri=$3, lokasi_kerja=$4,
		    gaji_range=$5, tanggal_mulai_kerja=$6, tanggal_selesai_kerja=$7,
		    status_pekerjaan=$8, deskripsi_pekerjaan=$9
		WHERE id=$10
	`
	_, err := database.DB.Exec(query,
		req.NamaPerusahaan, req.PosisiJabatan, req.BidangIndustri, req.LokasiKerja,
		req.GajiRange, req.TanggalMulaiKerja, tglSelesaiVal, req.StatusPekerjaan,
		req.DeskripsiPekerjaan, id,
	)
	if err != nil {
		return nil, err
	}
	return GetPekerjaanByIDRepo(id)
}

func GetPekerjaanRepo(search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error) {
	query := fmt.Sprintf(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja,
		       gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan,
		       deskripsi_pekerjaan, created_at, updated_at, is_delete, delete_by, deleted_at
		FROM pekerjaan_alumni
		WHERE nama_perusahaan ILIKE $1 
		   OR posisi_jabatan ILIKE $1 
		   OR bidang_industri ILIKE $1 
		   OR lokasi_kerja ILIKE $1
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortBy, order)

	rows, err := database.DB.Query(query, "%"+search+"%", limit, offset)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}
	defer rows.Close()

	var pekerjaanList []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		if err := rows.Scan(
			&p.ID,
			&p.AlumniID,
			&p.NamaPerusahaan,
			&p.PosisiJabatan,
			&p.BidangIndustri,
			&p.LokasiKerja,
			&p.GajiRange,
			&p.TanggalMulaiKerja,
			&p.TanggalSelesaiKerja,
			&p.StatusPekerjaan,
			&p.Deskripsi,
			&p.UpdatedAt,
			&p.CreatedAt,
			&p.IsDelete,
			&p.DeletedBy,
			&p.DeletedAt,
		); err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}

	return pekerjaanList, nil
}

func CountPekerjaanRepo(search string) (int, error) {
	var total int
	countQuery := `
		SELECT COUNT(*)
		FROM pekerjaan_alumni
		WHERE nama_perusahaan ILIKE $1 
		   OR posisi_jabatan ILIKE $1 
		   OR bidang_industri ILIKE $1 
		   OR lokasi_kerja ILIKE $1
	`
	err := database.DB.QueryRow(countQuery, "%"+search+"%").Scan(&total)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return total, nil
}

// Di dalam file: repository/pekerjaanRepo.go

func SoftDeletePekerjaan(pekerjaanID int, userID int, role string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer tx.Rollback()

	var currentOwnerUserID int
	var isDeleted sql.NullBool
	checkQuery := `
		SELECT a.user_id, pa.is_delete
		FROM pekerjaan_alumni pa
		JOIN alumni a ON pa.alumni_id = a.id
		WHERE pa.id = $1
	`
	err = tx.QueryRow(checkQuery, pekerjaanID).Scan(&currentOwnerUserID, &isDeleted)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrPekerjaanNotFound
		}
		return fmt.Errorf("error saat cek data pekerjaan: %w", err)
	}

	if isDeleted.Valid && isDeleted.Bool {
		return ErrPekerjaanNotFound
	}

	if role == "user" && currentOwnerUserID != userID {
		return ErrForbidden
	}

	queryUpdate := `
        UPDATE pekerjaan_alumni
        SET is_delete = TRUE,
            delete_by = $1,
            deleted_at = $2
        WHERE id = $3
	`
	_, err = tx.Exec(queryUpdate, userID, time.Now(), pekerjaanID)
	if err != nil {
		return fmt.Errorf("gagal melakukan soft delete: %w", err)
	}

	return tx.Commit()
}

func RestorePekerjaan(pekerjaanID int, userID int, role string) error {
	tx, err := database.DB.Begin()
	if err != nil {
		return fmt.Errorf("gagal memulai transaksi: %w", err)
	}
	defer tx.Rollback()

	if role == "user" {
		var ownerUserID int
		queryCekKepemilikan := `
			SELECT a.user_id
			FROM pekerjaan_alumni pa
			JOIN alumni a ON pa.alumni_id = a.id
			WHERE pa.id = $1 AND pa.is_delete = TRUE
		`
		err := tx.QueryRow(queryCekKepemilikan, pekerjaanID).Scan(&ownerUserID)
		if err != nil {
			if err == sql.ErrNoRows {
				return ErrPekerjaanNotFound
			}
			return fmt.Errorf("error saat cek kepemilikan: %w", err)
		}

		if ownerUserID != userID {
			return ErrForbidden
		}
	}

	queryUpdate := `
        UPDATE pekerjaan_alumni
        SET is_delete = FALSE,
            delete_by = NULL,
            deleted_at = NULL
        WHERE id = $1 AND COALESCE(is_delete, FALSE) = TRUE
	`
	result, err := tx.Exec(queryUpdate, pekerjaanID)
	if err != nil {
		return fmt.Errorf("gagal melakukan restore: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal memeriksa baris yang terpengaruh: %w", err)
	}

	if rowsAffected == 0 {
		return ErrPekerjaanNotFound
	}

	return tx.Commit()
}

func GetAllTrashPekerjaanRepo() ([]model.Pekerjaan, error) {
	query := `
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
		       lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
		       deskripsi_pekerjaan, is_delete, delete_by, deleted_at
		FROM pekerjaan_alumni
		WHERE COALESCE(is_delete, FALSE) = TRUE
		ORDER BY deleted_at DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil trash pekerjaan: %w", err)
	}
	defer rows.Close()

	var pekerjaanList []model.Pekerjaan
	for rows.Next() {
		var p model.Pekerjaan
		err := rows.Scan(
			&p.ID,
			&p.AlumniID,
			&p.NamaPerusahaan,
			&p.PosisiJabatan,
			&p.BidangIndustri,
			&p.LokasiKerja,
			&p.GajiRange,
			&p.TanggalMulaiKerja,
			&p.TanggalSelesaiKerja,
			&p.Deskripsi,
			&p.IsDelete,
			&p.DeletedBy,
			&p.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}

	return pekerjaanList, nil
}

func GetTrashPekerjaanByID(pekerjaanID int, userID int, role string) (*model.Pekerjaan, error) {
	var p model.Pekerjaan
	var err error

	// Kolom yang akan di-SELECT, sesuaikan dengan struct model.Pekerjaan Anda
	queryFields := `
		pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, pa.bidang_industri,
		pa.lokasi_kerja, pa.gaji_range, pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja,
		pa.deskripsi_pekerjaan, pa.is_delete, pa.delete_by, pa.deleted_at, pa.created_at, pa.updated_at
	`
	// Kumpulan variabel tujuan untuk hasil Scan
	scanDest := []interface{}{
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan,
		&p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja,
		&p.TanggalSelesaiKerja, &p.Deskripsi, &p.IsDelete,
		&p.DeletedBy, &p.DeletedAt, &p.CreatedAt, &p.UpdatedAt,
	}

	// Logika pemilihan query berdasarkan role
	if role == "admin" {
		// Admin bisa melihat trash manapun hanya berdasarkan ID pekerjaan
		query := fmt.Sprintf(`
			SELECT %s
			FROM pekerjaan_alumni pa
			WHERE pa.id = $1 AND pa.is_delete = TRUE
		`, queryFields)
		err = database.DB.QueryRow(query, pekerjaanID).Scan(scanDest...)
	} else {
		// User harus divalidasi kepemilikannya melalui join ke tabel alumni
		query := fmt.Sprintf(`
			SELECT %s
			FROM pekerjaan_alumni pa
			JOIN alumni a ON pa.alumni_id = a.id
			WHERE pa.id = $1 AND a.user_id = $2 AND pa.is_delete = TRUE
		`, queryFields)
		err = database.DB.QueryRow(query, pekerjaanID, userID).Scan(scanDest...)
	}

	// Menangani error, termasuk jika data tidak ditemukan (sql.ErrNoRows)
	if err != nil {
		return nil, err
	}

	return &p, nil
}


// Di dalam file: repository/pekerjaanRepo.go

func HardDeletePekerjaan(pekerjaanID int, userID int, role string) error {
	var query string
	var result sql.Result
	var err error

	if role == "admin" {

		query = `DELETE FROM pekerjaan_alumni WHERE id = $1 AND is_delete = TRUE`
		result, err = database.DB.Exec(query, pekerjaanID)
	} else {

		query = `
			DELETE FROM pekerjaan_alumni
			WHERE id = $1 AND is_delete = TRUE
			AND alumni_id IN (SELECT id FROM alumni WHERE user_id = $2)
		`
		result, err = database.DB.Exec(query, pekerjaanID, userID)
	}

	if err != nil {
		return fmt.Errorf("gagal mengeksekusi hard delete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal memeriksa baris yang terpengaruh: %w", err)
	}

	// Jika tidak ada baris terhapus, berarti data tidak ditemukan, belum di-soft delete, atau bukan milik user.
	if rowsAffected == 0 {
		return ErrPekerjaanNotFound
	}

	return nil
}