package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config interface {
	MongoUri() string
	MongoDatabaseName() string
	GrpcPort() int
}

type config struct {
	grpcPort      int
	mongoUser     string
	mongoPass     string
	mongoHost     string
	mongoPort     int
	mongoDatabase string
	mongoUri      string
}

func NewConfig() Config {
	var cfg config

	var err error

	cfg.grpcPort, err = strconv.Atoi(os.Getenv("GRPC_PORT"))

	if err != nil {
		log.Fatalln("Error parsing gRPC port env var \"GRPC_PORT\": ", err.Error())
	}

	value, exists := os.LookupEnv("MONGO_URI")
	if exists {
		cfg.mongoUri = value
		cfg.mongoDatabase = os.Getenv("MONGO_DATABASE_NAME")
		return &cfg
	}

	cfg.mongoUser = os.Getenv("MONGO_USERNAME")
	cfg.mongoPass = os.Getenv("MONGO_PASSWORD")
	cfg.mongoHost = os.Getenv("MONGO_HOST")
	cfg.mongoDatabase = os.Getenv("MONGO_DATABASE_NAME")

	cfg.mongoPort, err = strconv.Atoi(os.Getenv("MONGO_PORT"))

	if err != nil {
		log.Fatalln("Error parsing mongoDB port env var \"MONGO_PORT\": ", err.Error())
	}

	cfg.mongoUri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", cfg.mongoUser, cfg.mongoPass, cfg.mongoHost, cfg.mongoPort, cfg.mongoDatabase)

	return &cfg
}

func (c *config) MongoUri() string {
	return c.mongoUri
}

func (c *config) MongoDatabaseName() string {
	return c.mongoDatabase
}

func (c *config) GrpcPort() int {
	return c.grpcPort
}
