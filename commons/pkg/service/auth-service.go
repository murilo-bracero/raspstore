package service

import "github.com/murilo-bracero/raspstore/auth-service/proto/v1/auth-service/pb"

type AuthService interface {
	Authenticate(token string) (authResponse *pb.AuthenticateResponse, err error)
}
