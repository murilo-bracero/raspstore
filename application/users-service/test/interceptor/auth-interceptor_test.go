package middleware_test

import (
	"context"
	"log"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"raspstore.github.io/users-service/db"
	interceptor "raspstore.github.io/users-service/interceptors"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func mockHandler(ctx context.Context, req interface{}) (interface{}, error) {
	return "success", nil
}

func TestInterceptorWhenRouteIsNotWhitelisted(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())

	mdwr := interceptor.NewAuthInterceptor(cfg)

	m := make(map[string]string)
	m["authorization"] = "tokenmock"

	md := metadata.New(m)

	ctx := metadata.NewIncomingContext(context.Background(), md)

	info := &grpc.UnaryServerInfo{FullMethod: "/pb.AuthService/SignUp"}

	mdwr.WithUnaryAuthentication(ctx, "req", info, mockHandler)
}

func TestInterceptorWhenRouteIsWhitelisted(t *testing.T) {
	cfg := db.NewConfig()
	conn, err := db.NewMongoConnection(context.Background(), cfg)
	assert.NoError(t, err)
	defer conn.Close(context.Background())
	mdwr := interceptor.NewAuthInterceptor(cfg)

	m := make(map[string]string)
	m["authorization"] = "tokenmock"

	md := metadata.New(m)

	ctx := metadata.NewIncomingContext(context.Background(), md)

	info := &grpc.UnaryServerInfo{FullMethod: "/pb.AuthService/Login"}

	mdwr.WithUnaryAuthentication(ctx, "req", info, mockHandler)
}
