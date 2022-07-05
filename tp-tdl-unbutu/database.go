package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseGuard struct {
	Db     *mongo.Database
	client *mongo.Client
}

func ConnectToMongoDb(uri string) DatabaseGuard {
	opts := options.Client()
	opts.SetConnectTimeout(1 * time.Second)
	client, err := mongo.Connect(context.TODO(), opts.ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	return DatabaseGuard{client: client, Db: client.Database("TpTdl")}
}

func (db *DatabaseGuard) Disconnect() {
	if err := db.client.Disconnect(context.TODO()); err != nil {
		log.Fatal("Could not disconnect from database.", err)
	}
}
