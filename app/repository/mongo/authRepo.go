package mongo

import (
	"context"
	"latihan2/app/model/mongo"
	"latihan2/helper" // Mengimpor helper/base.go
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func GetUserByUsername(username string) (*mongo.User, error) {
	
	collection := helper.GetCollection("user")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user mongo.User
	filter := bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": username},
		},
		// "deleted_at": bson.M{"$exists": false}, // User yang di-soft-delete tidak bisa login
	}

	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}