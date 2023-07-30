package infra

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	MongoUri      string
	Database      string
	GrpcPort      int
	RestPort      int
	TokenSecret   string
	TokenDuration int
}

func NewConfig() *Config {
	return &Config{
		MongoUri:      os.Getenv("MONGO_URI"),
		Database:      os.Getenv("MONGO_DATABASE_NAME"),
		GrpcPort:      getIntEnv("GRPC_PORT"),
		RestPort:      getIntEnv("REST_PORT"),
		TokenSecret:   os.Getenv("JWT_SECRET"),
		TokenDuration: getIntEnv("TOKEN_DURATION"),
	}
}

func getIntEnv(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))

	if err != nil {
		log.Fatalf("error parsing env var %s: %s", key, err.Error())
	}

	return value
}
