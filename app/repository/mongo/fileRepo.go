package mongo

import (
	"context"
	"latihan2/app/model/mongo"
	"latihan2/helper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive" 

	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

func CreateFile(file *mongo.File) error {
	collection := helper.GetCollection("files") 
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	file.UploadedAt = time.Now()

	result, err := collection.InsertOne(ctx, file)
	if err != nil {
		return err
	}

	file.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func FindAllFiles() ([]mongo.File, error) {
	collection := helper.GetCollection("files")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var files []mongo.File
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &files); err != nil {
		return nil, err
	}
	return files, nil
}

func FindFileByID(id string) (*mongo.File, error) {
	collection := helper.GetCollection("files")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := helper.ToObjectID(id)
	if err != nil {
		return nil, err
	}

	var file mongo.File
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func OpenFileByID(id primitive.ObjectID) (*mongo.File, error) {
	collection := helper.GetCollection("files")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var file mongo.File
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&file)
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func DeleteFile(id string) error {
	collection := helper.GetCollection("files")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := helper.ToObjectID(id)
	if err != nil {
		return err
	}

	result, err := collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongodriver.ErrNoDocuments
	}
	return nil
}