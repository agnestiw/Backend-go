package mongo

import (
	"context"
	"latihan2/app/model/mongo"
	"latihan2/helper" // <-- Mengimpor helper/base.go
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetUsersRepo(search, sortBy, order string, limit, offset int) ([]mongo.User, error) {
    collection := helper.GetCollection("user") 
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var users []mongo.User

    filter := bson.M{
        "$or": []bson.M{
            {"username": bson.M{"$regex": search, "$options": "i"}},
            {"email": bson.M{"$regex": search, "$options": "i"}},
        },
        "deleted_at": nil,
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

    if err = cursor.All(ctx, &users); err != nil {
        return nil, err
    }
    return users, nil
}

func CountUsersRepo(search string) (int, error) {
	collection := helper.GetCollection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"username": bson.M{"$regex": search, "$options": "i"}},
			{"email": bson.M{"$regex": search, "$options": "i"}},
		},
		"deleted_at": nil,
	}

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// GetUserByID dipanggil oleh service.GetUsersByID
func GetUserByID(id string) (*mongo.User, error) {
	collection := helper.GetCollection("user") // <-- Menggunakan helper
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := helper.ToObjectID(id) // <-- Menggunakan helper
	if err != nil {
		return nil, err
	}

	var user mongo.User
	filter := bson.M{"_id": objID, "deleted_at": nil}
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}