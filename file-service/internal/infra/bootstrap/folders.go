package bootstrap

import (
	"context"
	"log/slog"
	"os"

	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type FolderBootstraper struct {
}

var _ Bootstraper = (*FolderBootstraper)(nil)

// Create all required folders for the application to run
func (b *FolderBootstraper) Bootstrap(ctx context.Context, config *config.Config) error {

	slog.Info("creating internal folder")
	if err := os.MkdirAll(config.Storage.Path+"/internal", os.ModePerm); err != nil {
		return err
	}

	slog.Info("creating secrets folder")
	if err := os.MkdirAll(config.Storage.Path+"/secrets", os.ModePerm); err != nil {
		return err
	}

	return nil
}
