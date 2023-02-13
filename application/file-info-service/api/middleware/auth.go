package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/murilo-bracero/raspstore-protofiles/auth-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"raspstore.github.io/file-manager/internal"
)

type userIdKey string

const UserIdKey userIdKey = "user-id-context-key"

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := r.Header.Get("Authorization")

		if token == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		conn, err := grpc.Dial(internal.AuthServiceUrl(), grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Println("[ERROR] Could not stablish connection to auth service, it may goes down: ", err.Error())
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		defer conn.Close()

		client := pb.NewAuthServiceClient(conn)

		in := &pb.AuthenticateRequest{Token: token}

		if res, err := client.Authenticate(context.Background(), in); err != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		} else {
			ctx := context.WithValue(r.Context(), UserIdKey, res.Uid)
			r = r.WithContext(ctx)
			log.Printf("[INFO] User %s is accessing resource %s", res.Uid, r.RequestURI)
		}

		h.ServeHTTP(w, r)
	})
}
