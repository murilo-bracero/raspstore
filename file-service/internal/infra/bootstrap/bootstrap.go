package bootstrap

import (
	"context"

	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type Bootstraper interface {
	Bootstrap(ctx context.Context, config *config.Config) error
}

var bootstrapers []Bootstraper = []Bootstraper{
	&FolderBootstraper{},
	&SecretsBootstraper{},
	&SQLiteBootstraper{},
}

func Run(ctx context.Context, config *config.Config) error {
	for _, bt := range bootstrapers {
		if err := bt.Bootstrap(ctx, config); err != nil {
			return err
		}
	}

	return nil
}
