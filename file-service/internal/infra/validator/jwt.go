package validator

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

type JWTValidator struct {
	config    *config.Config
	ar        *jwk.AutoRefresh
	publicKey *jwk.Key
}

const tokenPrefix = "Bearer"

var (
	ErrInvalidToken = errors.New("token is missing or is invalid")
)

func NewJWTValidator(ctx context.Context, config *config.Config) (*JWTValidator, error) {
	ar := jwk.NewAutoRefresh(ctx)

	ar.Configure(config.Auth.PublicKeyURL, jwk.WithMinRefreshInterval(15*time.Minute))

	if _, err := ar.Refresh(ctx, config.Auth.PublicKeyURL); err != nil {
		slog.Error("Failed to refresh JWT Tokens", "err", err)

		if !config.Auth.PAMEnabled {
			return nil, err
		}
	}

	publicKey, err := readPublicKey(config)

	if config.Auth.PAMEnabled && err != nil {
		return nil, err
	}

	return &JWTValidator{ar: ar, config: config, publicKey: publicKey}, nil
}

func readPublicKey(c *config.Config) (*jwk.Key, error) {
	pkFilePath := path.Join(c.Storage.Path, "secrets", "local-jwk.json")

	fpk, err := os.ReadFile(pkFilePath)

	if err != nil {
		slog.Error("Could not load Private Key to authenticate local PAM users", "err", err)
		return nil, err
	}

	jpk, err := jwk.ParseKey(fpk)

	if err != nil {
		slog.Error("Could not parse Private Key", "err", err)
		return nil, err
	}

	jPublicKey, err := jpk.PublicKey()
	jPublicKey.Set(jwk.KeyIDKey, "rstore")
	jPublicKey.Set(jwk.AlgorithmKey, jwa.RS256)

	if err != nil {
		slog.Error("Could not generate Public Key from Private Key", "err", err)
		return nil, err
	}

	return &jPublicKey, nil
}

func (p *JWTValidator) Validate(ctx context.Context, bearer string) (*jwt.Token, error) {
	token, err := getToken(bearer)

	if err != nil {
		return nil, err
	}

	keyset, err := p.ar.Fetch(ctx, p.config.Auth.PublicKeyURL)
	if err != nil {
		if !p.config.Auth.PAMEnabled {
			return nil, err
		}

		keyset = jwk.NewSet()
	}

	keyset.Add(*p.publicKey)

	jwt, err := jwt.Parse([]byte(token), jwt.WithKeySet(keyset), jwt.WithValidate(true))

	return &jwt, err
}

func getToken(bearer string) (token string, err error) {
	split := strings.Split(bearer, " ")

	if len(split) != 2 {
		return "", ErrInvalidToken
	}

	prefix := split[0]

	if prefix != tokenPrefix {
		return "", ErrInvalidToken
	}

	return split[1], nil
}
