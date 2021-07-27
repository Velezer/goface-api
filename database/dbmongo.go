package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	_ "go.mongodb.org/mongo-driver/mongo/readpref"
)

func client(ctx context.Context) *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panicln(err)
		}
	}()

	// err = client.Ping(ctx, readpref.Primary())
	// if err != nil {
	// 	log.Panicln(err)
	// }
	return client
}

func insertOne(data Face) *mongo.InsertOneResult {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := client(ctx)

	db := client.Database("krefa")
	coll := db.Collection("face")

	res, err := coll.InsertOne(ctx, data)
	if err != nil {
		log.Fatalln(err)
	}
	return res
}

func find() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client := client(ctx)

	db := client.Database("krefa")
	coll := db.Collection("face")

	cur, err := coll.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	var results []bson.M

	err = cur.All(ctx, &results)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range results {
		fmt.Println(v["descriptors"])

	}

}
