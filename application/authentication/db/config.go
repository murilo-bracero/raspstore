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
	GcpProjectId() string
	UserDataStorage() string
	CredentialsStorage() string
	TokenSecret() string
	TokenDuration() int
}

type config struct {
	grpcPort           int
	restPort           int
	gcpProjectId       string
	userDataStorage    string
	credentialsStorage string
	mongoUser          string
	mongoPass          string
	mongoHost          string
	mongoPort          int
	mongoDatabase      string
	mongoUri           string
	tokenSecret        string
	tokenDuration      int
}

func NewConfig() Config {
	var cfg config

	var err error

	cfg.userDataStorage = os.Getenv("USER_DATA_STORAGE")
	cfg.userDataStorage = os.Getenv("CREDENTIALS_STORAGE")
	cfg.gcpProjectId = os.Getenv("GCLOUD_PROJECT_ID")

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

	if value, err := strconv.Atoi(os.Getenv("TOKEN_DURATION")); err != nil {
		log.Fatalln("error parsing token duration env var TOKEN_DURATION: ", err.Error())
	} else {
		cfg.tokenDuration = value
	}

	cfg.tokenSecret = os.Getenv("JWT_SECRET")

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

func (c *config) GcpProjectId() string {
	return c.gcpProjectId
}

func (c *config) UserDataStorage() string {
	return c.userDataStorage
}

func (c *config) CredentialsStorage() string {
	return c.credentialsStorage
}

func (c *config) TokenSecret() string {
	return c.tokenSecret
}

func (c *config) TokenDuration() int {
	return c.tokenDuration
}
