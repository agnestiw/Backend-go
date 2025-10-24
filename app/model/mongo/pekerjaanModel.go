package mongo

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pekerjaan struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AlumniID            primitive.ObjectID `bson:"alumni_id,omitempty" json:"alumni_id"`
	NamaPerusahaan      string             `bson:"nama_perusahaan,omitempty" json:"nama_perusahaan"`
	PosisiJabatan       string             `bson:"posisi_jabatan,omitempty" json:"posisi_jabatan"`
	BidangIndustri      string             `bson:"bidang_industri,omitempty" json:"bidang_industri"`
	LokasiKerja         string             `bson:"lokasi_kerja,omitempty" json:"lokasi_kerja"`
	GajiRange           string             `bson:"gaji_range,omitempty" json:"gaji_range"`
	TanggalMulaiKerja   time.Time          `bson:"tanggal_mulai_kerja,omitempty" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time         `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string             `bson:"status_pekerjaan,omitempty" json:"status_pekerjaan"`
	Deskripsi           string             `bson:"deskripsi_pekerjaan,omitempty" json:"deskripsi_pekerjaan"`
	IsDelete            bool               `bson:"is_delete,omitempty" json:"is_delete"`
	DeleteBy            string             `bson:"delete_by,omitempty" json:"delete_by,omitempty"`
	DeletedAt           *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	CreatedAt           time.Time          `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt           *time.Time         `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type CreatePekerjaanRequest struct {
	ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AlumniID            primitive.ObjectID `bson:"alumni_id" json:"alumni_id"`
	NamaPerusahaan      string             `bson:"nama_perusahaan" json:"nama_perusahaan"`
	PosisiJabatan       string             `bson:"posisi_jabatan" json:"posisi_jabatan"`
	BidangIndustri      string             `bson:"bidang_industri" json:"bidang_industri"`
	LokasiKerja         string             `bson:"lokasi_kerja" json:"lokasi_kerja"`
	GajiRange           string             `bson:"gaji_range" json:"gaji_range"`
	TanggalMulaiKerja   time.Time          `bson:"tanggal_mulai_kerja" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time         `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string             `bson:"status_pekerjaan" json:"status_pekerjaan"`
	DeskripsiPekerjaan  string             `bson:"deskripsi_pekerjaan" json:"deskripsi_pekerjaan"`
	CreatedAt           time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt           time.Time          `bson:"updated_at" json:"updated_at"`
}

type UpdatePekerjaanRequest struct {
	NamaPerusahaan      string     `bson:"nama_perusahaan,omitempty" json:"nama_perusahaan,omitempty"`
	PosisiJabatan       string     `bson:"posisi_jabatan,omitempty" json:"posisi_jabatan,omitempty"`
	BidangIndustri      string     `bson:"bidang_industri,omitempty" json:"bidang_industri,omitempty"`
	LokasiKerja         string     `bson:"lokasi_kerja,omitempty" json:"lokasi_kerja,omitempty"`
	GajiRange           string     `bson:"gaji_range,omitempty" json:"gaji_range,omitempty"`
	TanggalMulaiKerja   *time.Time `bson:"tanggal_mulai_kerja,omitempty" json:"tanggal_mulai_kerja,omitempty"`
	TanggalSelesaiKerja *time.Time `bson:"tanggal_selesai_kerja,omitempty" json:"tanggal_selesai_kerja,omitempty"`
	StatusPekerjaan     string     `bson:"status_pekerjaan,omitempty" json:"status_pekerjaan,omitempty"`
	DeskripsiPekerjaan  string     `bson:"deskripsi_pekerjaan,omitempty" json:"deskripsi_pekerjaan,omitempty"`
	UpdatedAt           time.Time  `bson:"updated_at" json:"updated_at"`
}