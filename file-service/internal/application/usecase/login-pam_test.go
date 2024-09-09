package usecase_test

import (
	"context"
	"os"
	"os/exec"
	"os/user"
	"path"
	"testing"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/bootstrap"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {

	u, err := user.Current()

	assert.NoError(t, err, "user.Current")

	if u.Uid != "0" {
		t.Skip("This test requires root user priviledges to run")
	}

	username, password, err := createUser()

	if err != nil {
		t.Skip("This test requires root user priviledges to run")
	}

	mc := &config.Config{Storage: struct {
		Path  string
		Limit string
	}{Path: os.TempDir()}, Auth: struct {
		PAMEnabled   bool   "yaml:\"enable-pam\""
		PublicKeyURL string "yaml:\"public-key-url\""
	}{
		PAMEnabled: true,
	}}

	err = os.MkdirAll(path.Join(mc.Storage.Path, "secrets"), os.ModePerm)

	assert.NoError(t, err, "os.MkdirAll")

	sbs := &bootstrap.SecretsBootstrap{}

	err = sbs.Bootstrap(context.Background(), mc)

	assert.NoError(t, err, "sbs.Bootstrap")

	t.Cleanup(func() {
		err := deleteUser(username)

		if err != nil {
			t.Log("Could not delete user: username=", username)
		}
	})

	t.Run("Should return a JWT token if user credentials match", func(t *testing.T) {
		uc := usecase.NewLoginPAMUseCase(mc)

		token, err := uc.Execute(username, password)

		assert.NoError(t, err, "uc.Execute")

		assert.NotEmpty(t, token)
	})
}

func createUser() (string, string, error) {
	username := uuid.NewString()
	password := uuid.NewString()

	cmd := exec.Command("useradd", "-p", password, username)

	_, err := cmd.CombinedOutput()

	if err != nil {
		return "", "", err
	}

	return username, password, nil
}

func deleteUser(username string) error {
	cmd := exec.Command("userdel", "-f", username)

	_, err := cmd.CombinedOutput()

	return err
}
