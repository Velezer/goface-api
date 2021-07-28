package database

import (
	"context"
	"log"

	"github.com/Kagami/go-face"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsExist(ctx context.Context, coll *mongo.Collection, id string) {
	cur, err := coll.Find(ctx, bson.M{"_id": id})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var results []interface{}

	err = cur.All(ctx, &results)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(results)
}

func InsertOne(ctx context.Context, coll *mongo.Collection, data Face) (*mongo.InsertOneResult, error) {
	res, err := coll.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func PushDescriptor(ctx context.Context, coll *mongo.Collection, id interface{}, descriptor face.Descriptor) (*mongo.UpdateResult, error) {
	res, err := coll.UpdateByID(ctx, id, bson.M{"$push": bson.M{"descriptors": descriptor}})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FindAll(ctx context.Context, coll *mongo.Collection) (dataResult []Face) {
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	if err = cursor.All(ctx, &dataResult); err != nil {
		log.Fatal(err)
	}
	return
}
