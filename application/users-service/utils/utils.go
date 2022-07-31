package utils

import (
	"encoding/json"
	"net/http"

	"raspstore.github.io/users-service/validators"
)

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

	for _, oe := range validators.GetErrorsList() {
		if err == oe {
			code = http.StatusBadRequest
			break
		}
	}

	return code
}
