package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseGuard struct {
	Db     *mongo.Database
	client *mongo.Client
}

func ConnectToMongoDb(uri string) DatabaseGuard {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
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
