package service

import "github.com/murilo-bracero/raspstore-protofiles/auth-service/pb"

type AuthService interface {
	Authenticate(token string) (authResponse *pb.AuthenticateResponse, err error)
}
