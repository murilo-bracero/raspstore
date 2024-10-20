package bootstrap

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"path"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type SecretsBootstraper struct{}

var _ Bootstraper = (*SecretsBootstraper)(nil)

const bitsize = 4096

func (b *SecretsBootstraper) Bootstrap(ctx context.Context, config *config.Config) error {
	secretsPath := path.Join(config.Storage.Path, "secrets")

	pkFile, err := os.Open(secretsPath + "/key.json")

	if err != nil && !os.IsNotExist(err) {
		slog.Error("Could not open JWK file", "path", secretsPath, "err", err)
		return err
	}

	if os.IsNotExist(err) {
		err = createPrivateKey(ctx, secretsPath)

		return err
	}

	content, err := io.ReadAll(pkFile)

	if err != nil {
		slog.Error("Could not read jwk file", "err", err)
		return err
	}

	_, err = jwk.ParseKey(content)

	return err
}

func createPrivateKey(ctx context.Context, secretsPath string) error {
	privkey := jwk.NewRSAPrivateKey()

	rsaPk, err := rsa.GenerateKey(rand.Reader, bitsize)

	if err != nil {
		slog.Error("Could not generate private key", "err", err)
		return err
	}

	if err := privkey.FromRaw(rsaPk); err != nil {
		slog.Error("Could not read generated RSA private key", "err", err)
		return err
	}

	mpk, err := privkey.AsMap(ctx)

	if err != nil {
		slog.Error("Could not map private key", "err", err)
	}

	spk, err := json.Marshal(mpk)

	if err != nil {
		slog.Error("Could not convert private key map to json", "err", err)
		return err
	}

	fi, err := os.Create(secretsPath + "/key.json")

	if err != nil {
		slog.Error("Could not create jwk in internal secrets repository", "path", secretsPath, "err", err)
		return err
	}

	defer fi.Close()

	_, err = fi.Write(spk)

	if err != nil {
		slog.Error("Could not write new jwk to file", "path", secretsPath, "err", err)
		return err
	}

	return nil
}
