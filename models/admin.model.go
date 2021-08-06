package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

const collectionAdmin string = "coll_admin"

type Admin struct {
	Username string `bson:"admin" validate:"required"`
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
