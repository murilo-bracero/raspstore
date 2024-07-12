package handler_test

import (
	"testing"

	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/handler"
	"github.com/stretchr/testify/assert"
)

func TestValidateUpdateFileRequest(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		req := &model.UpdateFileRequest{
			Filename: "example.txt",
			Viewers:  []string{"user1", "user2"},
			Editors:  []string{"user1"},
		}

		err := handler.ValidateUpdateFileRequest(req)

		assert.NoError(t, err)
	})

	t.Run("should return error ErrFilenameEmpty", func(t *testing.T) {
		req := &model.UpdateFileRequest{
			Filename: "",
			Viewers:  []string{"user1", "user2"},
			Editors:  []string{"user1"},
		}

		err := handler.ValidateUpdateFileRequest(req)

		if err != handler.ErrFilenameEmpty {
			t.Errorf("Expected ErrFilenameEmpty, but got: %v", err)
		}
	})

	t.Run("should return error ErrViewersNil", func(t *testing.T) {
		req := &model.UpdateFileRequest{
			Filename: "example.txt",
			Viewers:  nil,
			Editors:  []string{"user1"},
		}

		err := handler.ValidateUpdateFileRequest(req)

		if err != handler.ErrViewersNil {
			t.Errorf("Expected ErrViewersNil, but got: %v", err)
		}
	})

	t.Run("should return error ErrEditorsNil", func(t *testing.T) {
		req := &model.UpdateFileRequest{
			Filename: "example.txt",
			Viewers:  []string{"user1", "user2"},
			Editors:  nil,
		}

		err := handler.ValidateUpdateFileRequest(req)

		if err != handler.ErrEditorsNil {
			t.Errorf("Expected ErrEditorsNil, but got: %v", err)
		}
	})
}
