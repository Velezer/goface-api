package database

import (
	"context"
	"goface-api/config"
	"goface-api/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DBRepo struct {
	RepoAdmin models.RepositoryAdminIface
	RepoFace  models.RepositoryFaceIface
}

func InitDB() *DBRepo {
	conf := config.GetDBConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.DB_URI))
	if err != nil {
		log.Panicln(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Panicln(err)
	}

	db := client.Database(conf.DB_NAME)
	dbrepo := DBRepo{
		RepoAdmin: &models.RepoAdmin{Collection: db.Collection("coll_admin")},
		RepoFace:  models.RepoFace{Collection: db.Collection("coll_face")},
	}
	return &dbrepo
}
