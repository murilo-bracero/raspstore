package db

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/datastore"
)

type DatastoreConnection interface {
	Client() *datastore.Client
	Close()
}

type datastoreConnection struct {
	client *datastore.Client
}

func NewDatastoreConnection(ctx context.Context, cfg Config) (DatastoreConnection, error) {

	if cfg.GcpProjectId() == "" {
		return nil, errors.New("env variable GCLOUD_PROJECT_ID is required when using datastore, but was not provided")
	}

	client, err := datastore.NewClient(ctx, cfg.GcpProjectId())

	if err != nil {
		return nil, err
	}

	return &datastoreConnection{client: client}, nil
}

func (d *datastoreConnection) Client() *datastore.Client {
	return d.client
}

func (d *datastoreConnection) Close() {
	err := d.client.Close()

	if err != nil {
		log.Fatalln("error while trying to close datastore client: ", err.Error())
	}
}
