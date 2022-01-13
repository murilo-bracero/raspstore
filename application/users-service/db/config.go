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
	RestPort() int
	AuthServiceUrl() string
}

type config struct {
	grpcPort       int
	restPort       int
	mongoUser      string
	mongoPass      string
	mongoHost      string
	mongoPort      int
	mongoDatabase  string
	mongoUri       string
	authServiceUrl string
}

func NewConfig() Config {
	var cfg config

	var err error

	if value, err := strconv.Atoi(os.Getenv("GRPC_PORT")); err != nil {
		log.Fatalln("error parsing gRPC port env var GRPC_PORT: ", err.Error())
	} else {
		cfg.grpcPort = value
	}

	if value, err := strconv.Atoi(os.Getenv("REST_PORT")); err != nil {
		log.Fatalln("error parsing gRPC port env var REST_PORT: ", err.Error())
	} else {
		cfg.restPort = value
	}

	cfg.authServiceUrl = os.Getenv("AUTH_SERVICE_URL")

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

func (c *config) RestPort() int {
	return c.restPort
}

func (c *config) AuthServiceUrl() string {
	return c.authServiceUrl
}
