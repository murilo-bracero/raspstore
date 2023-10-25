package validators_test

import (
	"testing"

	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model/request"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/validators"
	"github.com/stretchr/testify/assert"
)

func TestValidateUpdateFileRequest(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		req := &request.UpdateFileRequest{
			Filename: "example.txt",
			Viewers:  []string{"user1", "user2"},
			Editors:  []string{"user1"},
		}

		err := validators.ValidateUpdateFileRequest(req)

		assert.NoError(t, err)
	})

	t.Run("should return error ErrFilenameEmpty", func(t *testing.T) {
		req := &request.UpdateFileRequest{
			Filename: "",
			Viewers:  []string{"user1", "user2"},
			Editors:  []string{"user1"},
		}

		err := validators.ValidateUpdateFileRequest(req)

		if err != validators.ErrFilenameEmpty {
			t.Errorf("Expected ErrFilenameEmpty, but got: %v", err)
		}
	})

	t.Run("should return error ErrViewersNil", func(t *testing.T) {
		req := &request.UpdateFileRequest{
			Filename: "example.txt",
			Viewers:  nil,
			Editors:  []string{"user1"},
		}

		err := validators.ValidateUpdateFileRequest(req)

		if err != validators.ErrViewersNil {
			t.Errorf("Expected ErrViewersNil, but got: %v", err)
		}
	})

	t.Run("should return error ErrEditorsNil", func(t *testing.T) {
		req := &request.UpdateFileRequest{
			Filename: "example.txt",
			Viewers:  []string{"user1", "user2"},
			Editors:  nil,
		}

		err := validators.ValidateUpdateFileRequest(req)

		if err != validators.ErrEditorsNil {
			t.Errorf("Expected ErrEditorsNil, but got: %v", err)
		}
	})
}
