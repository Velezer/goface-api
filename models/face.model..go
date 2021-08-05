package models

import (
	"context"

	"github.com/Kagami/go-face"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


type Face struct {
	Id          string            `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string            `json:"name,omitempty" bson:"name,omitempty"`
	Descriptors []face.Descriptor `json:"descriptors,omitempty" bson:"descriptors,omitempty"`
}



func (face Face) InsertOne(ctx context.Context, coll *mongo.Collection, data Face) (*mongo.InsertOneResult, error) {
	res, err := coll.InsertOne(ctx, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (face Face) PushDescriptor(ctx context.Context, coll *mongo.Collection, id interface{}, descriptor face.Descriptor) (*mongo.UpdateResult, error) {
	res, err := coll.UpdateByID(ctx, id, bson.M{"$push": bson.M{"descriptors": descriptor}})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (face Face) FindAll(ctx context.Context, coll *mongo.Collection) (res []Face, err error) {
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &res); err != nil {
		return nil, err
	}
	return res, nil
}