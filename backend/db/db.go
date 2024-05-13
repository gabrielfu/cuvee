package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnector struct {
	URI        string
	Database   string
	Collection string
}

func NewMongoConnector(uri, database, collection string) *MongoConnector {
	return &MongoConnector{
		URI:        uri,
		Database:   database,
		Collection: collection,
	}
}

func (c *MongoConnector) Connect(ctx context.Context) (*mongo.Collection, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(c.URI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	db := client.Database(c.Database)
	collection := db.Collection(c.Collection)
	return collection, nil
}
