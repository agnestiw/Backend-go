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

func GetPekerjaanRepo(search, sortBy, order string, limit, offset int, role string, userID int) ([]model.Pekerjaan, error) {
    
    // 1. Buat query dasar
    baseQuery := `
        FROM pekerjaan_alumni pa
        LEFT JOIN alumni a ON pa.alumni_id = a.id
        WHERE (
            pa.nama_perusahaan ILIKE $1 
            OR pa.posisi_jabatan ILIKE $1 
            OR pa.bidang_industri ILIKE $1 
            OR pa.lokasi_kerja ILIKE $1
        ) AND pa.is_delete = false
    `
    // 2. Siapkan argumen
    args := []interface{}{"%" + search + "%"}
    
    // 3. Tambahkan filter user jika role bukan admin
    if role == "user" {
        baseQuery += " AND a.user_id = $2"
        args = append(args, userID)
    }

    // 4. Buat query SELECT
    selectQuery := fmt.Sprintf(`
        SELECT pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, pa.bidang_industri, pa.lokasi_kerja,
               pa.gaji_range, pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja, pa.status_pekerjaan,
               pa.deskripsi_pekerjaan, pa.created_at, pa.updated_at
        %s
        ORDER BY %s %s
        LIMIT $%d OFFSET $%d
    `, baseQuery, sortBy, order, len(args)+1, len(args)+2) // Penomoran $ sudah benar
    
    args = append(args, limit, offset)

    // 5. Eksekusi
    rows, err := database.DB.Query(selectQuery, args...)
    if err != nil {
        log.Println("Query error:", err)
        return nil, err
    }
    defer rows.Close()

    // 6. Scan hasilnya (INI BAGIAN YANG PERLU DILENGKAPI)
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
            &p.CreatedAt,
            &p.UpdatedAt,
        ); err != nil {
            return nil, err // Tangani error saat scanning
        }
        pekerjaanList = append(pekerjaanList, p)
    }

    // Selalu cek error setelah loop rows
    if err = rows.Err(); err != nil {
        return nil, err
    }

    return pekerjaanList, nil
}

func CountPekerjaanRepo(search string, role string, userID int) (int, error) {
    var total int
    baseQuery := `
        FROM pekerjaan_alumni pa
        LEFT JOIN alumni a ON pa.alumni_id = a.id
        WHERE (
            pa.nama_perusahaan ILIKE $1 
            OR pa.posisi_jabatan ILIKE $1 
            OR pa.bidang_industri ILIKE $1 
            OR pa.lokasi_kerja ILIKE $1
        ) AND pa.is_delete = false
    `
    args := []interface{}{"%" + search + "%"}

    if role == "user" {
        baseQuery += " AND a.user_id = $2"
        args = append(args, userID)
    }
    countQuery := fmt.Sprintf("SELECT count(*) %s", baseQuery)

    err := database.DB.QueryRow(countQuery, args...).Scan(&total)
    if err != nil {
        log.Printf("Error counting pekerjaan: %v", err)
        return 0, err
    }

    return total, nil
}

func GetPekerjaanByIDRepo(id int, userID int, role string) (*model.Pekerjaan, error) {
	var p model.Pekerjaan
	
    // Kolom-kolom yang akan di-SELECT
	queryFields := `
        pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, pa.bidang_industri,
        pa.lokasi_kerja, pa.gaji_range, pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja,
        pa.status_pekerjaan, pa.deskripsi_pekerjaan, pa.created_at, pa.updated_at
    ` // <-- 'updated_at' sudah ditambahkan
    
    // Tujuan Scan
	scanDest := []interface{}{
		&p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan,
		&p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja,
		&p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.Deskripsi, &p.CreatedAt,
		&p.UpdatedAt, // <-- 'UpdatedAt' sudah ditambahkan
	}

	var query string
    var err error

	if role == "admin" {
        // Admin bisa lihat data manapun (yang belum di-soft-delete)
		query = fmt.Sprintf(`
            SELECT %s
            FROM pekerjaan_alumni pa
            WHERE pa.id = $1 AND pa.is_delete = false
        `, queryFields)
		err = database.DB.QueryRow(query, id).Scan(scanDest...)
	
    } else {
        // User hanya bisa lihat data miliknya sendiri (yang belum di-soft-delete)
		query = fmt.Sprintf(`
            SELECT %s
            FROM pekerjaan_alumni pa
            JOIN alumni a ON pa.alumni_id = a.id
            WHERE pa.id = $1 AND a.user_id = $2 AND pa.is_delete = false
        `, queryFields)
		err = database.DB.QueryRow(query, id, userID).Scan(scanDest...)
	}

	if err != nil {
		return nil, err
	}
	return &p, nil
}

func GetPekerjaanByAlumniID(alumniID int) ([]model.Pekerjaan, error) {
	rows, err := database.DB.Query(`
		SELECT id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
			   lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
			   status_pekerjaan, deskripsi_pekerjaan, created_at
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
			&p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.Deskripsi, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		pekerjaanList = append(pekerjaanList, p)
	}
	return pekerjaanList, nil
}

func CreatePekerjaan(req model.CreatePekerjaanRequest) (*model.Pekerjaan, error) {
	// Gunakan model.Pekerjaan (dari prompt pertama) sebagai target Scan
	var p model.Pekerjaan

	// Ini adalah query SQL LENGKAP tanpa '...'
	query := `
    INSERT INTO pekerjaan_alumni 
    (
        alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
        lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
        status_pekerjaan, deskripsi_pekerjaan, is_delete, created_at
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, FALSE, $11)
    RETURNING 
        id, alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri,
        lokasi_kerja, gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja,
        status_pekerjaan, deskripsi_pekerjaan, created_at, updated_at
    `

	// --- Konversi Tipe Data (Sama seperti sebelumnya) ---
	tglMulai, err := time.Parse("2006-01-02", req.TanggalMulaiKerja)
	if err != nil {
		return nil, fmt.Errorf("format tanggal_mulai_kerja salah: %w", err)
	}

	var tglSelesaiVal sql.NullTime
	if req.TanggalSelesaiKerja != "" {
		if t, err := time.Parse("2006-01-02", req.TanggalSelesaiKerja); err == nil {
			tglSelesaiVal = sql.NullTime{Time: t, Valid: true}
		} else {
			return nil, fmt.Errorf("format tanggal_selesai_kerja salah: %w", err)
		}
	}
	// --- Selesai Konversi ---

	// Eksekusi query dengan argumen yang benar
	err = database.DB.QueryRow(query,
		req.AlumniID,           // $1
		req.NamaPerusahaan,     // $2
		req.PosisiJabatan,      // $3
		req.BidangIndustri,     // $4
		req.LokasiKerja,        // $5
		req.GajiRange,          // $6
		tglMulai,               // $7 (tipe time.Time)
		tglSelesaiVal,          // $8 (tipe sql.NullTime)
		req.StatusPekerjaan,    // $9
		req.DeskripsiPekerjaan, // $10
		time.Now(),             // $11 (untuk created_at)
	).Scan(
		// Scan ke struct model.Pekerjaan (p)
		// Urutan ini HARUS cocok dengan urutan RETURNING
		&p.ID,
		&p.AlumniID,
		&p.NamaPerusahaan,
		&p.PosisiJabatan,
		&p.BidangIndustri,
		&p.LokasiKerja,
		&p.GajiRange,
		&p.TanggalMulaiKerja,   // Target adalah time.Time
		&p.TanggalSelesaiKerja, // Target adalah sql.NullTime
		&p.StatusPekerjaan,
		&p.Deskripsi, // Target adalah string (deskripsi_pekerjaan -> Deskripsi)
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		log.Println("SQL Insert/Scan Error:", err)
		return nil, err
	}

	// Kembalikan struct Pekerjaan yang sudah terisi
	return &p, nil
}

func UpdatePekerjaan(id int, userID int, role string, req model.UpdatePekerjaanRequest) (*model.Pekerjaan, error) {
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
	return GetPekerjaanByIDRepo(id, userID, role)
}

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

func GetTrashPekerjaanByID(pekerjaanID int, userID int, role string) (*model.TrashPekerjaanResponse, error) {
    var p model.TrashPekerjaanResponse
    var err error

    // PERBAIKAN: Tambahkan 3 kolom soft delete di akhir
    queryFields := `
        pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, pa.bidang_industri,
        pa.lokasi_kerja, pa.gaji_range, pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja,
        pa.status_pekerjaan, pa.deskripsi_pekerjaan, pa.created_at, pa.updated_at,
        pa.is_delete, pa.delete_by, pa.deleted_at
    `

    // PERBAIKAN: Tambahkan 3 field tujuan untuk soft delete di akhir
    scanDest := []interface{}{
        &p.ID, &p.AlumniID, &p.NamaPerusahaan, &p.PosisiJabatan,
        &p.BidangIndustri, &p.LokasiKerja, &p.GajiRange, &p.TanggalMulaiKerja,
        &p.TanggalSelesaiKerja, &p.StatusPekerjaan, &p.Deskripsi,
        &p.CreatedAt, &p.UpdatedAt,
        &p.IsDelete, &p.DeletedBy, &p.DeletedAt, // <-- Ditambahkan di sini
    }

    if role == "admin" {
        query := fmt.Sprintf(`
            SELECT %s
            FROM pekerjaan_alumni pa
            WHERE pa.id = $1 AND pa.is_delete = TRUE
        `, queryFields)
        // Gunakan ... untuk "membongkar" slice scanDest
        err = database.DB.QueryRow(query, pekerjaanID).Scan(scanDest...) 
    } else {
        query := fmt.Sprintf(`
            SELECT %s
            FROM pekerjaan_alumni pa
            JOIN alumni a ON pa.alumni_id = a.id
            WHERE pa.id = $1 AND a.user_id = $2 AND pa.is_delete = TRUE
        `, queryFields)
        // Gunakan ... untuk "membongkar" slice scanDest
        err = database.DB.QueryRow(query, pekerjaanID, userID).Scan(scanDest...) 
    }

    if err != nil {
        return nil, err
    }

    return &p, nil
}

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

	if rowsAffected == 0 {
		return ErrPekerjaanNotFound
	}

	return nil
}
