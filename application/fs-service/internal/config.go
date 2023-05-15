package internal

import (
	"log"
	"os"
	"strconv"
)

func FileInfoServiceUrl() string {
	return os.Getenv("FILE_INFO_SERVICE_URL")
}

func AuthServiceUrl() string {
	return os.Getenv("AUTH_SERVICE_URL")
}

func RestPort() int {
	return getIntEnv("REST_PORT")
}

func StoragePath() string {
	return os.Getenv("STORAGE_PATH")
}

func getIntEnv(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))

	if err != nil {
		log.Fatalf("error parsing env var %s: %s", key, err.Error())
	}

	return value
}
