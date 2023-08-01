package database

import (
	"context"
	"fmt"

	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/file-info-service/internal"
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
	fmt.Println("Connecting to MongoDB...")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(internal.MongoUri()))

	if err != nil {
		logger.Error("Could not connect to MongoDB: %s", err.Error())
		return nil, err
	}

	fmt.Println("Connected to MongoDB Successfully")
	return &conn{database: client.Database(internal.MongoDatabaseName())}, nil
}

func (c *conn) Close(ctx context.Context) {
	err := c.database.Client().Disconnect(ctx)

	if err != nil {
		logger.Error("Error releasing MongoDB connection: %s", err.Error())
	}
}

func (c *conn) Collection(collectionName string) *mongo.Collection {
	return c.database.Collection(collectionName)
}
