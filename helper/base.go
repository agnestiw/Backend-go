package helper

import (
	"errors"
	"latihan2/database" // Sesuaikan path ke file database Anda

	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

// GetCollection diawali huruf kapital (Exported)
// Fungsi ini mengambil koleksi dari variabel database.MongoDB Anda
func GetCollection(colName string) *mongodriver.Collection {
	return database.MongoDB.Collection(colName)
}

// ToObjectID diawali huruf kapital (Exported)
func ToObjectID(id string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return objID, errors.New("format ID tidak valid")
	}
	return objID, nil
}