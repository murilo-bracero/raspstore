package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"raspstore.github.io/file-manager/api/dto"
	"raspstore.github.io/file-manager/db"
	"raspstore.github.io/file-manager/pb"
)

type AuthMiddleware interface {
	Apply(h http.Handler) http.Handler
}

type authMiddleware struct {
	cfg db.Config
}

func NewAuthMiddleware(cfg db.Config) AuthMiddleware {
	return &authMiddleware{cfg: cfg}
}

func (a *authMiddleware) Apply(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			er := new(dto.ErrorResponse)
			er.Message = "authorization header is missing"
			er.Code = "AM01"
			send(w, er)
			return
		}

		conn, err := grpc.Dial(a.cfg.AuthServiceUrl())

		if err != nil {
			log.Fatalln("could not stablish connection to auth service, it may goes down: ", err.Error())
		}

		defer conn.Close()

		client := pb.NewAuthServiceClient(conn)

		in := &pb.AuthenticateRequest{Token: token}

		if res, err := client.Authenticate(context.Background(), in); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			er := new(dto.ErrorResponse)
			er.Message = "authorization header is missing"
			er.Reason = err.Error()
			er.Code = "AM01"
			send(w, er)
			return
		} else {
			log.Printf("user %s is accessing resource %s", res.Uid, r.RequestURI)
		}

		h.ServeHTTP(w, r)
	})
}

func send(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse, err := json.Marshal(obj)
	if err != nil {
		return
	}
	w.Write(jsonResponse)
}
