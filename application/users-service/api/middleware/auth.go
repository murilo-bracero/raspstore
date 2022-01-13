package middleware

import (
	"context"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"raspstore.github.io/users-service/api/dto"
	"raspstore.github.io/users-service/db"
	"raspstore.github.io/users-service/pb"
	"raspstore.github.io/users-service/utils"
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
			utils.Send(w, dto.ErrorResponse{Message: "authorization header is missing", Code: "AM01"})
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
			utils.Send(w, dto.ErrorResponse{Message: "authorization header is missing", Reason: err.Error(), Code: "AM02"})
			return
		} else {
			r.Header.Add("UID", res.Uid)
			log.Printf("user %s is accessing resource %s", res.Uid, r.RequestURI)
		}

		h.ServeHTTP(w, r)
	})
}
