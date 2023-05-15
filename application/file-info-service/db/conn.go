package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"raspstore.github.io/file-manager/internal"
)

type MongoConnection interface {
	Close(ctx context.Context)
	DB() *mongo.Database
}

type conn struct {
	database *mongo.Database
}

func NewMongoConnection(ctx context.Context) (MongoConnection, error) {
	fmt.Println("Connecting to MongoDB...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(internal.MongoUri()))

	if err != nil {
		log.Fatalln("Could not connect to MongoDB: ", err.Error())
		return nil, err
	}

	fmt.Println("Connected to MongoDB Successfully")
	return &conn{database: client.Database(internal.MongoDatabaseName())}, nil
}

func (c *conn) Close(ctx context.Context) {
	err := c.database.Client().Disconnect(ctx)

	if err != nil {
		log.Fatalln("Error releasing MongoDB connection: ", err.Error())
	}
}

func (c *conn) DB() *mongo.Database {
	return c.database
}
