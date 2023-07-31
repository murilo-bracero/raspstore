package database

import (
	"context"
	"fmt"
	"log"

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
	fmt.Println("Connecting to MongoDB...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoUri))

	if err != nil {
		log.Fatalln("Could not connect to MongoDB: ", err.Error())
		return nil, err
	}

	fmt.Println("Connected to MongoDB Successfully")
	return &conn{database: client.Database(config.Database)}, nil
}

func (c *conn) Close(ctx context.Context) {
	err := c.database.Client().Disconnect(ctx)

	if err != nil {
		log.Fatalln("Error releasing MongoDB connection: ", err.Error())
	}
}

func (c *conn) Collection(collectionName string) *mongo.Collection {
	return c.database.Collection(collectionName)
}
