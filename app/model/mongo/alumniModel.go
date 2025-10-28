package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Alumni memetakan data di collection 'alumni'
type Alumni struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"user_id" bson:"user_id"` // Foreign key ke collection 'users'
	NIM        string             `json:"nim" bson:"nim"`
	Nama       string             `json:"nama" bson:"nama"`
	Jurusan    string             `json:"jurusan" bson:"jurusan"`
	Angkatan   int                `json:"angkatan" bson:"angkatan"`
	TahunLulus int                `json:"tahun_lulus" bson:"tahun_lulus"`
	Email      string             `json:"email" bson:"email"`
	NoTelepon  *string            `json:"no_telepon,omitempty" bson:"no_telepon,omitempty"`
	Alamat     *string            `json:"alamat,omitempty" bson:"alamat,omitempty"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt  *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
