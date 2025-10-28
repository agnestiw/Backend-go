package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User memetakan data di collection 'users'
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Email     string             `json:"email" bson:"email"`
	Password  string             `json:"-" bson:"password"` // Penting: Jangan tampilkan password di JSON
	Role      string             `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}