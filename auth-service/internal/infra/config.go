package infra

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	MongoUri          string
	Database          string
	GrpcPort          int
	RestPort          int
	TokenSecret       string
	TokenDuration     int
	MinPasswordLength int
	EnforceMfa        bool
}

func NewConfig() *Config {
	return &Config{
		MongoUri:          os.Getenv("MONGO_URI"),
		Database:          os.Getenv("MONGO_DATABASE_NAME"),
		GrpcPort:          getIntEnv("GRPC_PORT"),
		RestPort:          getIntEnv("REST_PORT"),
		TokenSecret:       os.Getenv("JWT_SECRET"),
		TokenDuration:     getIntEnv("TOKEN_DURATION"),
		MinPasswordLength: getIntEnv("MIN_PASS_LEN"),
		EnforceMfa:        getBoolEnv("ENFORCE_MFA"),
	}
}

func getIntEnv(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))

	if err != nil {
		log.Fatalf("error parsing env var %s: %s", key, err.Error())
	}

	return value
}

func getBoolEnv(key string) bool {
	value, err := strconv.ParseBool(os.Getenv(key))

	if err != nil {
		return false
	}

	return value
}
