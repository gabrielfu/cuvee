package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnector struct {
	URI      string
	Database string
}

func NewMongoConnector(uri, database string) *MongoConnector {
	return &MongoConnector{
		URI:      uri,
		Database: database,
	}
}

func (c *MongoConnector) Connect(ctx context.Context) (*mongo.Database, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(c.URI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	db := client.Database(c.Database)
	return db, nil
}
