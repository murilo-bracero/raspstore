package usecase_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	m "github.com/murilo-bracero/raspstore/file-service/internal/infra/middleware"
	"github.com/stretchr/testify/assert"
)

type mockFilesRepository struct{}

func (m *mockFilesRepository) FindUsageByUserId(userID string) (int64, error) {
	return 0, errors.New("Not implemented!")
}

func (m *mockFilesRepository) Save(file *entity.File) error {
	return errors.New("Not implemented!")
}

func (m *mockFilesRepository) Delete(userId string, fileId string) error {
	if strings.HasSuffix(fileId, "_error") {
		return errors.New("generic error")
	}
	return nil
}

func (m *mockFilesRepository) FindAll(userId string, page int, size int, filename string, secret bool) (filesPage *entity.FilePage, err error) {
	return nil, errors.New("Not implemented!")
}

func (m *mockFilesRepository) FindById(userId string, fileId string) (*entity.File, error) {
	return nil, errors.New("Not implemented!")
}

func (m *mockFilesRepository) Update(userId string, file *entity.File) error {
	return errors.New("Not implemented!")
}

func TestDeleteFileUseCase_Execute(t *testing.T) {
	useCase := usecase.NewDeleteFileUseCase(&mockFilesRepository{})

	token := jwt.New()
	token.Set("sub", uuid.NewString())

	ctx := context.WithValue(context.WithValue(context.Background(),
		middleware.RequestIDKey, "trace12345"),
		m.UserClaimsCtxKey, token)

	t.Run("ValidFileDeletion", func(t *testing.T) {
		fileId := "file123"

		err := useCase.Execute(ctx, fileId)

		assert.NoError(t, err)
	})

	t.Run("FileDeletionError", func(t *testing.T) {
		useCase := usecase.NewDeleteFileUseCase(&mockFilesRepository{})

		fileId := "file456_error"

		err := useCase.Execute(ctx, fileId)

		assert.Error(t, err)
	})
}
