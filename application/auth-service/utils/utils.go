package utils

import (
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/grpc/metadata"
)

type contextKey string

var (
	ContextKeyUID = contextKey("uid")
)

func GetValueFromMetadata(key string, ctx context.Context) string {
	md, exists := metadata.FromIncomingContext(ctx)

	if !exists {
		return ""
	}

	values := md[key]

	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func Send(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(obj)
	if err != nil {
		return
	}
	w.Write(jsonResponse)
}

func ReqStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	code := http.StatusInternalServerError

	for _, oe := range GetErrorsList() {
		if err == oe {
			code = http.StatusBadRequest
			break
		}
	}

	return code
}
