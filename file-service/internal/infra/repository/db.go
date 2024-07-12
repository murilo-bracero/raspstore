package repository

import (
	"context"
	"log/slog"

	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseConnection interface {
	Close(ctx context.Context)
	Collection(collectionName string) *mongo.Collection
}

type databaseConnection struct {
	database *mongo.Database
}

func NewDatabaseConnection(ctx context.Context, config *config.Config) (*databaseConnection, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.Database.Uri))

	if err != nil {
		slog.Error("Could not connect to database", "error", err)
		return nil, err
	}

	return &databaseConnection{database: client.Database(config.Database.Name)}, nil
}

func (c *databaseConnection) Close(ctx context.Context) {
	err := c.database.Client().Disconnect(ctx)

	if err != nil {
		slog.Error("Error releasing database connection", "error", err)
	}
}

func (c *databaseConnection) Collection(collectionName string) *mongo.Collection {
	return c.database.Collection(collectionName)
}
