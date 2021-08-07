package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionAdmin string = "coll_admin"

type Admin struct {
	Username string `bson:"_id" validate:"required"`
	Password string `bson:"password" validate:"required"`
}

func (admin Admin) InsertOne(ctx context.Context, db *mongo.Database) (*mongo.InsertOneResult, error) {
	coll := db.Collection(collectionAdmin)
	res, err := coll.InsertOne(ctx, admin)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (admin Admin) FindOneByID(ctx context.Context, db *mongo.Database) (res Admin, err error) {
	coll := db.Collection(collectionAdmin)
	cursor := coll.FindOne(ctx, bson.M{"_id": admin.Username})

	if err = cursor.Decode(&res); err != nil {
		return res, err
	}
	return res, nil
}
