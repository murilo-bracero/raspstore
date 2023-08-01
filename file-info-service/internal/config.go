package internal

import (
	"fmt"
	"os"
	"strconv"

	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
)

func GrpcPort() int {
	return getIntEnv("GRPC_PORT")
}

func MongoUri() string {
	if value, exists := os.LookupEnv("MONGO_URI"); exists {
		return value
	}

	mongoUser := os.Getenv("MONGO_USERNAME")
	mongoPass := os.Getenv("MONGO_PASSWORD")
	mongoHost := os.Getenv("MONGO_HOST")

	mongoDatabase := MongoDatabaseName()

	mongoPort := getIntEnv("MONGO_PORT")

	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", mongoUser, mongoPass, mongoHost, mongoPort, mongoDatabase)
}

func MongoDatabaseName() string {
	return os.Getenv("MONGO_DATABASE_NAME")
}

func AuthServiceUrl() string {
	return os.Getenv("AUTH_SERVICE_URL")
}

func UserServiceUrl() string {
	return os.Getenv("USER_SERVICE_URL")
}

func RestPort() int {
	return getIntEnv("REST_PORT")
}

func StorageLimit() string {
	return os.Getenv("STORAGE_LIMIT")
}

func StoragePath() string {
	return os.Getenv("STORAGE_PATH")
}

func getIntEnv(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))

	if err != nil {
		logger.Error("error parsing env var %s: %s", key, err.Error())
		os.Exit(1)
	}

	return value
}
