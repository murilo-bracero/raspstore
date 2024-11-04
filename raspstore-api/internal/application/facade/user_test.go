package facade_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/bootstrap"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserSave(t *testing.T) {

	config := &config.Config{Storage: config.StorageConfig{Path: os.TempDir()}}

	err := (&bootstrap.FolderBootstraper{}).Bootstrap(context.Background(), config)

	assert.NoError(t, err, "FolderBootstraper")

	err = (&bootstrap.SecretsBootstraper{}).Bootstrap(context.Background(), config)

	assert.NoError(t, err, "SecretsBootstraper")

	t.Run("should save user", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		mockObj := mocks.NewMockRepository[entity.User](mockCtrl)
		mockObj.EXPECT().Save(gomock.Any()).Return(nil)

		ff := facade.NewUserFacade(config, mockObj)

		user := &model.CreateUserRequest{
			Username: "user1",
			Password: "password1",
		}

		err := ff.Save(user)

		assert.NoError(t, err)
	})

	t.Run("should return error if error on repository", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)

		mockObj := mocks.NewMockRepository[entity.User](mockCtrl)
		mockObj.EXPECT().Save(gomock.Any()).Return(errors.New("error on repository"))

		ff := facade.NewUserFacade(config, mockObj)

		user := &model.CreateUserRequest{
			Username: "user1",
			Password: "password1",
		}

		err := ff.Save(user)

		assert.Error(t, err)
	})

	t.Cleanup(func() {
		os.RemoveAll(os.TempDir() + "/secrets")
	})
}
