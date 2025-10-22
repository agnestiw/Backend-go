package model

import "time"

// Alumni (API Model) - Struct generik yang dikembalikan ke klien
// Perhatikan: ID dan UserID diubah menjadi string agar konsisten.
// sql.NullString diubah menjadi *string (pointer)
type Alumni struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	NIM        string    `json:"nim"`
	Nama       string    `json:"nama"`
	Jurusan    string    `json:"jurusan"`
	Angkatan   int       `json:"angkatan"`
	TahunLulus int       `json:"tahun_lulus"`
	Email      string    `json:"email"`
	NoTelepon  *string   `json:"no_telepon,omitempty"`
	Alamat     *string   `json:"alamat,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// CreateAlumniRequest (API Model)
// Perhatikan: UserID diubah menjadi string
type CreateAlumniRequest struct {
	NIM        string `json:"nim"`
	UserID     string `json:"user_id"`
	Nama       string `json:"nama"`
	Jurusan    string `json:"jurusan"`
	Angkatan   int    `json:"angkatan"`
	TahunLulus int    `json:"tahun_lulus"`
	Email      string `json:"email"`
	NoTelepon  string `json:"no_telepon"`
	Alamat     string `json:"alamat"`
}

// UpdateAlumniRequest (API Model)
// Struct ini tidak perlu diubah karena tidak mengandung ID
type UpdateAlumniRequest struct {
	Nama       string `json:"nama"`
	Jurusan    string `json:"jurusan"`
	Angkatan   int    `json:"angkatan"`
	TahunLulus int    `json:"tahun_lulus"`
	Email      string `json:"email"`
	NoTelepon  string `json:"no_telepon"`
	Alamat     string `json:"alamat"`
}

// AlumniPekerjaanResponse (API Model)
// Perhatikan: ID diubah menjadi string
type AlumniPekerjaanResponse struct {
	ID             string `json:"id"`
	Jurusan        string `json:"jurusan"`
	TahunLulus     int    `json:"tahun_lulus"`
	BidangIndustri string `json:"bidang_industri"`
	NamaPerusahaan string `json:"nama_perusahaan"`
	PosisiJabatan  string `json:"posisi_jabatan"`
	GajiRange      string `json:"gaji_range"`
}