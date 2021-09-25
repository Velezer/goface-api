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

type RepoFace struct {
	Collection *mongo.Collection
}

func (repo RepoFace) InsertOne(face Face) error {
	_, err := repo.Collection.InsertOne(context.Background(), face)
	if err != nil {
		return err
	}
	return nil
}

func (repo RepoFace) PushDescriptor(id string, descriptor face.Descriptor) (*mongo.UpdateResult, error) {
	res, err := repo.Collection.UpdateByID(context.Background(), id, bson.M{"$push": bson.M{"descriptors": descriptor}})
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo RepoFace) FindById(id string) (res []Face, err error) {
	cursor, err := repo.Collection.Find(context.Background(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (repo RepoFace) FindAll() (res []Face, err error) {
	cursor, err := repo.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.Background(), &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (repo RepoFace) DeleteId(id string) (*mongo.DeleteResult, error) {
	res, err := repo.Collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return nil, err
	}
	return res, nil
}
