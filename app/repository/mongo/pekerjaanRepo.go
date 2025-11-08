package mongo

import (
	"context"
	"errors"
	"fmt"
	model "latihan2/app/model/mongo"
	"latihan2/database"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var pekerjaanColl *mongo.Collection
var (
	ErrPekerjaanNotFound = errors.New("pekerjaan not found")
	ErrForbidden         = errors.New("forbidden access")
)

func InitPekerjaanCollection(db *mongo.Database) {
	pekerjaanColl = db.Collection("pekerjaan")
}

func GetPekerjaanRepo(search, sortBy, order string, limit, offset int) ([]model.Pekerjaan, error) {
	pekerjaanColl := database.MongoDB.Collection("pekerjaan")
	if pekerjaanColl == nil {
		return nil, errors.New("pekerjaanColl belum diinisialisasi")
	}
	filter := bson.M{
		"$or": []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
			{"lokasi_kerja": bson.M{"$regex": search, "$options": "i"}},
		},
	}
	sortOrder := 1
	if order == "desc" {
		sortOrder = -1
	}
	opts := options.Find().
		SetSort(bson.D{{Key: sortBy, Value: sortOrder}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))
	cursor, err := pekerjaanColl.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var pekerjaanList []model.Pekerjaan
	if err := cursor.All(context.TODO(), &pekerjaanList); err != nil {
		return nil, err
	}
	return pekerjaanList, nil
}

func CountPekerjaanRepo(search string) (int, error) {
	if pekerjaanColl == nil {
		return 0, errors.New("pekerjaanColl belum diinisialisasi")
	}
	filter := bson.M{
		"$or": []bson.M{
			{"nama_perusahaan": bson.M{"$regex": search, "$options": "i"}},
			{"posisi_jabatan": bson.M{"$regex": search, "$options": "i"}},
			{"bidang_industri": bson.M{"$regex": search, "$options": "i"}},
			{"lokasi_kerja": bson.M{"$regex": search, "$options": "i"}},
		},
	}
	count, err := pekerjaanColl.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func GetPekerjaanByIDRepo(id string) (*model.Pekerjaan, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var p model.Pekerjaan
	err = pekerjaanColl.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&p)
	if err != nil {
		return nil, ErrPekerjaanNotFound
	}
	return &p, nil
}

func GetPekerjaanByAlumniID(alumniID string) ([]model.Pekerjaan, error) {
	objID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		// fmt.Println("DEBUG: alumni_id bukan ObjectID valid:", alumniID)
		return nil, err
	}

	filter := bson.M{"alumni_id": objID}
	cursor, err := pekerjaanColl.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	var pekerjaanList []model.Pekerjaan
	if err := cursor.All(context.TODO(), &pekerjaanList); err != nil {
		return nil, err
	}
	fmt.Printf("DEBUG: Total data ditemukan: %d\n", len(pekerjaanList))
	return pekerjaanList, nil
}

func CreatePekerjaan(p *model.Pekerjaan) (*model.Pekerjaan, error) {
	p.ID = primitive.NewObjectID()
	p.CreatedAt = time.Now()
	p.IsDelete = false
	if _, err := pekerjaanColl.InsertOne(context.TODO(), p); err != nil {
		return nil, err
	}
	return p, nil
}

func UpdatePekerjaan(id string, req model.UpdatePekerjaanRequest) (*model.Pekerjaan, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	update := bson.M{
		"$set": bson.M{
			"nama_perusahaan":     req.NamaPerusahaan,
			"posisi_jabatan":      req.PosisiJabatan,
			"bidang_industri":     req.BidangIndustri,
			"lokasi_kerja":        req.LokasiKerja,
			"gaji_range":          req.GajiRange,
			"tanggal_mulai_kerja": req.TanggalMulaiKerja,
			"status_pekerjaan":    req.StatusPekerjaan,
			"deskripsi_pekerjaan": req.DeskripsiPekerjaan,
			"updated_at":          time.Now(),
		},
	}
	_, err = pekerjaanColl.UpdateByID(context.TODO(), objID, update)
	if err != nil {
		return nil, err
	}
	return GetPekerjaanByIDRepo(id)
}

func SoftDeletePekerjaan(pekerjaanID, userID, role string) error {
	objID, err := primitive.ObjectIDFromHex(pekerjaanID)
	if err != nil {
		return err
	}
	var pekerjaan model.Pekerjaan
	err = pekerjaanColl.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&pekerjaan)
	if err != nil {
		return ErrPekerjaanNotFound
	}
	if pekerjaan.IsDelete {
		return ErrPekerjaanNotFound
	}
	update := bson.M{
		"$set": bson.M{
			"is_delete":  true,
			"delete_by":  userID,
			"deleted_at": time.Now(),
		},
	}
	result, err := pekerjaanColl.UpdateByID(context.TODO(), objID, update)
	if err != nil {
		fmt.Println("DEBUG: Gagal update:", err)
		return err
	}
	fmt.Printf("DEBUG: Update result: %+v\n", result)
	return nil
}

func RestorePekerjaan(pekerjaanID, userID, role string) error {
	objID, err := primitive.ObjectIDFromHex(pekerjaanID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID, "is_delete": true}
	if role == "user" {
		filter["user_id"] = userID
	}

	update := bson.M{
		"$set": bson.M{
			"is_delete":  false,
			"delete_by":  nil,
			"deleted_at": nil,
		},
	}

	result, err := pekerjaanColl.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrPekerjaanNotFound
	}
	return nil
}

func HardDeletePekerjaan(pekerjaanID, userID, role string) error {
	objID, err := primitive.ObjectIDFromHex(pekerjaanID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID, "is_delete": true}
	if role == "user" {
		filter["user_id"] = userID
	}

	result, err := pekerjaanColl.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return ErrPekerjaanNotFound
	}
	return nil
}

func GetTrashPekerjaan(id, role string) ([]model.Pekerjaan, error) {
	filter := bson.M{"is_delete": true}

	// Jika ID dikirim lewat params, ambil berdasarkan _id
	if id != "" {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, fmt.Errorf("invalid ID format: %v", err)
		}
		filter["_id"] = objID
	}

	// Jika role adalah "user", tolak akses karena trash hanya untuk admin
	if role == "user" {
		return nil, fmt.Errorf("forbidden: user cannot access trash")
	}

	fmt.Println("DEBUG FILTER:", filter) // untuk bantu debug

	cursor, err := pekerjaanColl.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var pekerjaanList []model.Pekerjaan
	if err := cursor.All(context.TODO(), &pekerjaanList); err != nil {
		return nil, err
	}

	if len(pekerjaanList) == 0 {
		return nil, fmt.Errorf("Pekerjaan not found")
	}

	return pekerjaanList, nil
}
