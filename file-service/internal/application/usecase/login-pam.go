package usecase

import (
	"errors"
	"log/slog"
	"os"
	"path"
	"time"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/msteinert/pam"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

var (
	ErrUnknowMessageStyle = errors.New("unknown pam message style")
	ErrPAMFailed          = errors.New("pam authentication failed")
	ErrNotAuthenticated   = errors.New("not authenticated")
	ErrPAMNotEnabled      = errors.New("pam not enable for this domain")
)

type LoginPAMUseCase interface {
	Execute(username, password string) (string, error)
}

type loginPAMUseCase struct {
	config *config.Config
}

var _ LoginPAMUseCase = (*loginPAMUseCase)(nil)

func NewLoginPAMUseCase(config *config.Config) *loginPAMUseCase {
	return &loginPAMUseCase{config: config}
}

func (l *loginPAMUseCase) Execute(username, password string) (string, error) {
	if !l.config.Auth.PAMEnabled {
		return "", ErrPAMNotEnabled
	}

	t, err := pam.StartFunc("passwd", username, func(s pam.Style, msg string) (string, error) {
		switch s {
		case pam.PromptEchoOff:
			return password, nil
		case pam.PromptEchoOn:
			slog.Warn("pam authentication with promp echo on without password", "username", username, "msg", msg)
			return "", nil
		case pam.ErrorMsg:
			slog.Warn("pam authentication failed", "username", username, "msg", msg)
			return "", nil
		case pam.TextInfo:
			slog.Info("pam authentication message info", "username", username, "msg", msg)
			return "", nil
		default:
			slog.Error("unknow pam message style", "username", username, "msg", msg)
			return "", ErrUnknowMessageStyle
		}
	})

	if err != nil {
		slog.Error("pam message trade failed", "username", username)
		return "", ErrPAMFailed
	}

	err = t.Authenticate(0)

	if err != nil {
		slog.Warn("auth attempt failed", "username", username)
		return "", ErrNotAuthenticated
	}

	privateKey, err := l.readPrivateKey()

	if err != nil {
		slog.Error("could not read Private Key", "err", err)
		return "", ErrNotAuthenticated
	}

	token, err := l.generateJWT(username, privateKey)

	if err != nil {
		slog.Error("could not generate JWT", "err", nil)
		return "", ErrNotAuthenticated
	}

	return token, nil
}

func (l *loginPAMUseCase) readPrivateKey() (*jwk.Key, error) {
	pkPath := path.Join(l.config.Storage.Path, "secrets", "key.json")

	fpk, err := os.ReadFile(pkPath)

	if err != nil {
		slog.Error("Could not read private key", "err", err)
		return nil, err
	}

	jPrivateKey, err := jwk.ParseKey(fpk)
	jPrivateKey.Set(jwk.KeyIDKey, "rstore")

	if err != nil {
		slog.Error("Could not parse private key", "err", err)
		return nil, err
	}

	return &jPrivateKey, nil
}

func (l *loginPAMUseCase) generateJWT(subject string, pk *jwk.Key) (string, error) {
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
