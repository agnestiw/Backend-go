package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func InitMongoDB() {
	mongoURI := os.Getenv("MONGODB_URI")
	dbName := os.Getenv("DATABASE_NAME")
	if mongoURI == "" {
		log.Fatal("Error: MONGODB_URI tidak ditemukan di file .env. Pastikan file .env sudah disimpan.")
	}
	if dbName == "" {
		log.Fatal("Error: DATABASE_NAME tidak ditemukan di file .env. Pastikan file .env sudah disimpan.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Gagal koneksi ke MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Tidak bisa ping ke MongoDB: %v", err)
	}

	MongoClient = client
	MongoDB = client.Database(dbName)

	fmt.Println("Berhasil tersambung ke MongoDB!")
	collections, err := MongoDB.ListCollectionNames(ctx, map[string]interface{}{})
	if err != nil {
		log.Fatalf("Gagal mendapatkan list koleksi: %v", err)
	}
	fmt.Println("Koleksi yang ada:", collections)

	log.Println("Berhasil tersambung ke MongoDB!")
}