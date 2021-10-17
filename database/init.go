package database

import (
	"context"
	"goface-api/iface"
	"goface-api/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type DBRepo struct {
	RepoAdmin iface.RepositoryAdminIface
	RepoFace  iface.RepositoryFaceIface
}

func InitDB(uri string, dbname string) (*DBRepo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	db := client.Database(dbname)
	dbrepo := DBRepo{
		RepoAdmin: models.RepoAdmin{Collection: db.Collection("coll_admin")},
		RepoFace:  models.RepoFace{Collection: db.Collection("coll_face")},
	}
	return &dbrepo, nil
}
