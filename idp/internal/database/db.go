package database

import (
	"context"
	"os"

	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/idp/internal/infra"
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
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoUri))

	if err != nil {
		logger.Error("Could not connect to MongoDB: %s", err.Error())
		return nil, err
	}

	logger.Info("Connected to MongoDB Successfully")
	return &conn{database: client.Database(config.Database)}, nil
}

func (c *conn) Close(ctx context.Context) {
	err := c.database.Client().Disconnect(ctx)

	if err != nil {
		logger.Error("Error releasing MongoDB connection: %s", err.Error())
		os.Exit(1)
	}
}

func (c *conn) Collection(collectionName string) *mongo.Collection {
	return c.database.Collection(collectionName)
}
