package database

import (
	"context"
	"log/slog"

	"github.com/murilo-bracero/raspstore/file-service/internal"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConnection interface {
	Close(ctx context.Context)
	Collection(collectionName string) *mongo.Collection
}

type conn struct {
	database *mongo.Database
}

func NewMongoConnection(ctx context.Context) (MongoConnection, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(internal.MongoUri()))

	if err != nil {
		slog.Error("Could not connect to MongoDB: %s", err.Error())
		return nil, err
	}

	return &conn{database: client.Database(internal.MongoDatabaseName())}, nil
}

func (c *conn) Close(ctx context.Context) {
	err := c.database.Client().Disconnect(ctx)

	if err != nil {
		slog.Error("Error releasing MongoDB connection: %s", err.Error())
	}
}

func (c *conn) Collection(collectionName string) *mongo.Collection {
	return c.database.Collection(collectionName)
}
