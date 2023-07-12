package grpc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"raspstore.github.io/users-service/internal/grpc"
	"raspstore.github.io/users-service/internal/model"
	"raspstore.github.io/users-service/proto/v1/users-service/pb"
)

func TestGetUserConfiguration(t *testing.T) {
	svc := &userConfigServiceMock{}
	subject := grpc.NewUserGrpcService(svc)

	result, err := subject.GetUserConfiguration(context.Background(), &pb.GetUserConfigurationRequest{})

	assert.NoError(t, err)

	assert.Equal(t, "2G", result.StorageLimit)
	assert.Equal(t, int64(8), result.MinPasswordLength)
	assert.Equal(t, true, result.AllowPublicUserCreation)
	assert.Equal(t, false, result.AllowLoginWithEmail)
	assert.Equal(t, false, result.EnforceMfa)
}

func TestGetUserConfigurationError(t *testing.T) {
	svc := &userConfigServiceMock{shouldThrowError: true}
	subject := grpc.NewUserGrpcService(svc)

	_, err := subject.GetUserConfiguration(context.Background(), &pb.GetUserConfigurationRequest{})

	assert.Error(t, err)
}

type userConfigServiceMock struct {
	shouldThrowError bool
}

func (u *userConfigServiceMock) UpdateUserConfig(userConfig *model.UserConfiguration) error {
	return nil
}

func (u *userConfigServiceMock) GetUserConfig() (userConfig *model.UserConfiguration, err error) {
	if u.shouldThrowError {
		return nil, errors.New("generic error")
	}

	return &model.UserConfiguration{
		StorageLimit:            "2G",
		MinPasswordLength:       8,
		AllowPublicUserCreation: true,
		AllowLoginWithEmail:     false,
		EnforceMfa:              false,
	}, nil
}

func (u *userConfigServiceMock) ValidateUser(user *model.User, isAdminCreation bool) error {
	return nil
}
