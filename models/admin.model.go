package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Admin struct {
	Username string `bson:"_id" form:"username" json:"username" validate:"required"`
	Password string `bson:"password" form:"password" json:"password" validate:"required"`
}


type RepoAdmin struct {
	Collection *mongo.Collection
}

func (repo *RepoAdmin) FindOneByID(id string) (res Admin, err error) {
	cursor := repo.Collection.FindOne(context.Background(), bson.M{"_id": id})

	if err = cursor.Decode(&res); err != nil {
		return res, err
	}
	return res, nil
}

func (repo *RepoAdmin) InsertOne(admin Admin) error {
	_, err := repo.Collection.InsertOne(context.Background(), admin)
	if err != nil {
		return err
	}
	return nil
}
