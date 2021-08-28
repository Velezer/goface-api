package models

import (
	"context"

	"github.com/Kagami/go-face"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const collectionFace string = "coll_terserah"

type Face struct {
	Id          string            `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string            `json:"name,omitempty" bson:"name,omitempty"`
	Descriptors []face.Descriptor `json:"descriptors,omitempty" bson:"descriptors,omitempty"`
}

func (face Face) InsertOne(ctx context.Context, db *mongo.Database) (*mongo.InsertOneResult, error) {
	coll := db.Collection(collectionFace)
	res, err := coll.InsertOne(ctx, face)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (face Face) PushDescriptor(ctx context.Context, db *mongo.Database) (*mongo.UpdateResult, error) {
	coll := db.Collection(collectionFace)
	res, err := coll.UpdateByID(ctx, face.Id, bson.M{"$push": bson.M{"descriptors": face.Descriptors[0]}})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (face Face) FindById(ctx context.Context, db *mongo.Database) (res []Face, err error) {
	coll := db.Collection(collectionFace)
	cursor, err := coll.Find(ctx, bson.M{"_id": face.Id})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (face Face) FindAll(ctx context.Context, db *mongo.Database) (res []Face, err error) {
	coll := db.Collection(collectionFace)
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (face Face) DeleteId(ctx context.Context, db *mongo.Database) (*mongo.DeleteResult, error) {
	coll := db.Collection(collectionFace)
	res, err := coll.DeleteOne(ctx, bson.M{"_id": face.Id})
	if err != nil {
		return nil, err
	}
	return res, nil
}
