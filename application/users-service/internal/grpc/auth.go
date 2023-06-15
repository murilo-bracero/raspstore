package grpc

import (
	"context"
	"log"

	"github.com/murilo-bracero/raspstore-protofiles/auth-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"raspstore.github.io/users-service/internal"
)

func Authenticate(token string) (uid string, err error) {
	conn, err := grpc.Dial(internal.AuthServiceUrl(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("[ERROR] Could not stablish connection to auth service :", err.Error())
		return "", err
	}

	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	in := &pb.AuthenticateRequest{Token: token}

	if res, err := client.Authenticate(context.Background(), in); err != nil {
		return "", err
	} else {
		return res.Uid, nil
	}
}
