package auth

import (
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

func GenerateRefreshToken(config *config.Config) (string, error) {
	pk, err := readPrivateKey(config.Storage.Path)

	if err != nil {
		slog.Error("Could not read private key", "err", err)
		return "", err
	}

	return generateJWT(uuid.NewString(), pk)
}

func readPrivateKey(storagePath string) (*jwk.Key, error) {
	pkPath := path.Join(storagePath, "secrets", "key.json")

	fpk, err := os.ReadFile(pkPath)

	if err != nil {
		slog.Error("Could not read private key", "err", err)
		return nil, err
	}

	jPrivateKey, err := jwk.ParseKey(fpk)
	if err := jPrivateKey.Set(jwk.KeyIDKey, "rstore"); err != nil {
		slog.Error("Could not set KeyId in private key", "err", err)
		return nil, err
	}

	if err != nil {
		slog.Error("Could not parse private key", "err", err)
		return nil, err
	}

	return &jPrivateKey, nil
}

func generateJWT(subject string, pk *jwk.Key) (string, error) {
	tkn, err := jwt.NewBuilder().
		Issuer("rstore").
		IssuedAt(time.Now()).
		Subject(subject).
		Expiration(time.Now().Add(1 * time.Hour)).
		Build()

	if err != nil {
		slog.Error("Could not create JWT", "err", err)
		return "", err
	}

	signed, err := jwt.Sign(tkn, jwa.RS256, *pk)

	if err != nil {
		slog.Error("Could not sign JWT", "err", err)
		return "", err
	}

	return string(signed), nil
}
