package db

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config interface {
	GrpcPort() int
	RootFolder() string
	MongoUri() string
	MongoDatabaseName() string
	AuthServiceUrl() string
	RestPort() int
}

type config struct {
	grpcPort       int
	rootFolder     string
	mongoUser      string
	mongoHost      string
	mongoDatabase  string
	mongoPass      string
	mongoPort      int
	mongoUri       string
	authServiceUrl string
	restPort       int
}

func NewConfig() Config {
	var cfg config

	if value, err := strconv.Atoi(os.Getenv("GRPC_PORT")); err != nil {
		log.Fatalln("error parsing gRPC port env var GRPC_PORT: ", err.Error())
	} else {
		cfg.grpcPort = value
	}

	if value, err := strconv.Atoi(os.Getenv("REST_PORT")); err != nil {
		log.Fatalln("error parsing rest api port env var REST_PORT: ", err.Error())
	} else {
		cfg.restPort = value
	}

	cfg.rootFolder = os.Getenv("ROOT_FOLDER")
	cfg.authServiceUrl = os.Getenv("AUTH_SERVICE_URL")

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

	if value, err := strconv.Atoi(os.Getenv("MONGO_PORT")); err != nil {
		log.Fatalln("Error parsing mongoDB port env var \"MONGO_PORT\": ", err.Error())
	} else {
		cfg.mongoPort = value
	}

	cfg.mongoUri = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", cfg.mongoUser, cfg.mongoPass, cfg.mongoHost, cfg.mongoPort, cfg.mongoDatabase)

	return &cfg
}

func (c *config) GrpcPort() int {
	return c.grpcPort
}

func (c *config) RootFolder() string {
	return c.rootFolder
}

func (c *config) MongoUri() string {
	return c.mongoUri
}

func (c *config) MongoDatabaseName() string {
	return c.mongoDatabase
}

func (c *config) AuthServiceUrl() string {
	return c.authServiceUrl
}

func (c *config) RestPort() int {
	return c.restPort
}