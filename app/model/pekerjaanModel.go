package model

import (
	"database/sql"
	"time"
)

type Pekerjaan struct {
    ID                 int            `json:"id"`
    AlumniID           int            `json:"alumni_id"`
    NamaPerusahaan     string         `json:"nama_perusahaan"`
    PosisiJabatan      string         `json:"posisi_jabatan"`
    BidangIndustri     string         `json:"bidang_industri"`
    LokasiKerja        string         `json:"lokasi_kerja"`
    GajiRange          string         `json:"gaji_range"`
    TanggalMulaiKerja  time.Time      `json:"tanggal_mulai_kerja"`
    TanggalSelesaiKerja sql.NullTime  `json:"tanggal_selesai_kerja"`
    StatusPekerjaan    string         `json:"status_pekerjaan"`
    Deskripsi          string         `json:"deskripsi_pekerjaan"`
    IsDelete           sql.NullBool   `json:"is_delete"`
    DeletedBy          sql.NullInt64  `json:"delete_by"`
    DeletedAt          sql.NullTime   `json:"deleted_at"`
    CreatedAt          time.Time      `json:"created_at"`
    UpdatedAt          time.Time      `json:"updated_at"`
}

type TrashPekerjaanResponse struct {
	ID                 int          `json:"id"`
	AlumniID           int          `json:"alumni_id"`
	NamaPerusahaan     string       `json:"nama_perusahaan"`
	PosisiJabatan      string       `json:"posisi_jabatan"`
	BidangIndustri     string       `json:"bidang_industri"`
	LokasiKerja        string       `json:"lokasi_kerja"`
	GajiRange          string       `json:"gaji_range"`
	TanggalMulaiKerja  time.Time    `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja sql.NullTime `json:"tanggal_selesai_kerja"`
	StatusPekerjaan    string       `json:"status_pekerjaan"`
	Deskripsi          string       `json:"deskripsi_pekerjaan"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
}

type CreatePekerjaanRequest struct {
	AlumniID            int    `json:"alumni_id"`
	NamaPerusahaan      string `json:"nama_perusahaan"`
	PosisiJabatan       string `json:"posisi_jabatan"`
	BidangIndustri      string `json:"bidang_industri"`
	LokasiKerja         string `json:"lokasi_kerja"`
	GajiRange           string `json:"gaji_range"`
	TanggalMulaiKerja   string `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja string `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string `json:"status_pekerjaan"`
	DeskripsiPekerjaan  string `json:"deskripsi_pekerjaan"`
}

type UpdatePekerjaanRequest struct {
	NamaPerusahaan      string `json:"nama_perusahaan"`
	PosisiJabatan       string `json:"posisi_jabatan"`
	BidangIndustri      string `json:"bidang_industri"`
	LokasiKerja         string `json:"lokasi_kerja"`
	GajiRange           string `json:"gaji_range"`
	TanggalMulaiKerja   string `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja string `json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string `json:"status_pekerjaan"`
	DeskripsiPekerjaan  string `json:"deskripsi_pekerjaan"`
}
