package database

import (
	"context"
	"log/slog"

	"github.com/murilo-bracero/raspstore/file-service/internal/infra"
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

func NewMongoConnection(ctx context.Context, config *infra.Config) (MongoConnection, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Database.Uri))

	if err != nil {
		slog.Error("Could not connect to MongoDB", "error", err)
		return nil, err
	}

	return &conn{database: client.Database(config.Database.Name)}, nil
}

func (c *conn) Close(ctx context.Context) {
	err := c.database.Client().Disconnect(ctx)

	if err != nil {
		slog.Error("Error releasing MongoDB connection", "error", err)
	}
}

func (c *conn) Collection(collectionName string) *mongo.Collection {
	return c.database.Collection(collectionName)
}
