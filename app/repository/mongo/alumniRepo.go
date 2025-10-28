package mongo

import (
	"context"
	"latihan2/app/model" // DTO Bersama (CreateAlumniRequest)
	"latihan2/app/model/mongo"
	"latihan2/helper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAlumniRepo(search, sortBy, order string, limit, offset int) ([]mongo.Alumni, error) {
	collection := helper.GetCollection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var alumniList []mongo.Alumni

	filter := bson.M{
		"$or": []bson.M{
			{"nim": bson.M{"$regex": search, "$options": "i"}},
			{"nama": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		},
		"deleted_at": bson.M{"$exists": false}, // Filter soft delete
	}

	orderVal := 1
	if order == "desc" {
		orderVal = -1
	}
	opts := options.Find().
		SetSort(bson.D{{Key: sortBy, Value: orderVal}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &alumniList); err != nil {
		return nil, err
	}
	return alumniList, nil
}

// CountAlumniRepo dipanggil oleh service.GetAllAlumni
func CountAlumniRepo(search string) (int, error) {
	collection := helper.GetCollection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"nim": bson.M{"$regex": search, "$options": "i"}},
			{"nama": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		},
		"deleted_at": bson.M{"$exists": false}, // Filter soft delete
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetAlumniByID dipanggil oleh service.GetAlumniByID
func GetAlumniByID(id string) (*mongo.Alumni, error) {
	collection := helper.GetCollection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := helper.ToObjectID(id)
	if err != nil {
		return nil, err
	}

	var a mongo.Alumni
	filter := bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}}
	err = collection.FindOne(ctx, filter).Decode(&a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// CreateAlumni dipanggil oleh service.CreateAlumni
func CreateAlumni(alumni *mongo.Alumni) (*mongo.Alumni, error) {
	collection := helper.GetCollection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	alumni.ID = primitive.NewObjectID()
	alumni.CreatedAt = time.Now()
	alumni.UpdatedAt = time.Now()

	_, err := collection.InsertOne(ctx, alumni)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

// UpdateAlumni dipanggil oleh service.UpdateAlumni
func UpdateAlumni(id string, req model.UpdateAlumniRequest) (*mongo.Alumni, error) {
	collection := helper.GetCollection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := helper.ToObjectID(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"nama":        req.Nama,
			"jurusan":     req.Jurusan,
			"angkatan":    req.Angkatan,
			"tahun_lulus": req.TahunLulus,
			"email":       req.Email,
			"no_telepon":  req.NoTelepon,
			"alamat":      req.Alamat,
			"updated_at":  time.Now(),
		},
	}

	filter := bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	if result.MatchedCount == 0 {
		return nil, mongodriver.ErrNoDocuments
	}

	return GetAlumniByID(id)
}

// SoftDeleteAlumni dipanggil oleh service.SoftDeleteAlumni
func SoftDeleteAlumni(id string) error {
	collection := helper.GetCollection("alumni")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := helper.ToObjectID(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"deleted_at": time.Now(),
		},
	}

	filter := bson.M{"_id": objID, "deleted_at": bson.M{"$exists": false}}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongodriver.ErrNoDocuments
	}
	return nil
}