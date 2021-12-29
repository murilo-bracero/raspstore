package middleware_test

import (
	"context"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"raspstore.github.io/authentication/db"
	interceptor "raspstore.github.io/authentication/interceptors"
	"raspstore.github.io/authentication/repository"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

type mockTokenManager struct{}

func (m *mockTokenManager) Generate(uid string) (string, error) {
	return "tokenmock", nil
}

func (m *mockTokenManager) Verify(token string) (string, error) {
	return "uidmock", nil
}

func mockHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return "success", nil
}

func TestInterceptorWhenRouteIsNotWhitelisted(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	userRepo := repository.NewUsersRepository(context.Background(), conn)
	tokenManager := new(mockTokenManager)

	mdwr := interceptor.NewAuthInterceptor(userRepo, tokenManager)

	m := make(map[string]string)
	m["authorization"] = "tokenmock"

	md := metadata.New(m)

	ctx := metadata.NewIncomingContext(context.Background(), md)

	info := &grpc.UnaryServerInfo{FullMethod: "/pb.AuthService/SignUp"}

	mdwr.WithAuthentication(ctx, "req", info, mockHandler)
}

func TestInterceptorWhenRouteIsWhitelisted(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	userRepo := repository.NewUsersRepository(context.Background(), conn)
	tokenManager := new(mockTokenManager)

	mdwr := interceptor.NewAuthInterceptor(userRepo, tokenManager)

	m := make(map[string]string)
	m["authorization"] = "tokenmock"

	md := metadata.New(m)

	ctx := metadata.NewIncomingContext(context.Background(), md)

	info := &grpc.UnaryServerInfo{FullMethod: "/pb.AuthService/Login"}

	mdwr.WithAuthentication(ctx, "req", info, mockHandler)
}
