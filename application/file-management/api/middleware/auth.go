package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/murilo-bracero/raspstore-protofiles/authentication/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"raspstore.github.io/file-manager/api/dto"
	"raspstore.github.io/file-manager/db"
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
			return
		}

		conn, err := grpc.Dial(a.cfg.AuthServiceUrl(), grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Println("could not stablish connection to auth service, it may goes down: ", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			send(w, dto.ErrorResponse{Message: "feature temporatily unavailable.", Code: "AM99"})
			return
		}

		defer conn.Close()

		client := pb.NewAuthServiceClient(conn)

		in := &pb.AuthenticateRequest{Token: token}

		if res, err := client.Authenticate(context.Background(), in); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		} else {
			r.Header.Add("UID", res.Uid)
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
