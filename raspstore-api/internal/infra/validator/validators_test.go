package validator_test

import (
	"testing"

	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/validator"
	"github.com/stretchr/testify/assert"
)

func TestValidateUpdateFileRequest(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		req := &model.UpdateFileRequest{
			Filename: "example.txt",
		}

		err := validator.ValidateUpdateFileRequest(req)

		assert.NoError(t, err)
	})

	t.Run("should return error ErrFilenameEmpty", func(t *testing.T) {
		req := &model.UpdateFileRequest{
			Filename: "",
		}

		err := validator.ValidateUpdateFileRequest(req)

		if err != validator.ErrFilenameEmpty {
			t.Errorf("Expected ErrFilenameEmpty, but got: %v", err)
		}
	})
}
