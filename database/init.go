package database

import (
	"context"
	"goface-api/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	DB *mongo.Database
)

func InitDB() {
	conf := config.GetDBConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.DB_URI))
	if err != nil {
		log.Panicln(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Panicln(err)
	}

	DB = client.Database(conf.DB_NAME)
	
}