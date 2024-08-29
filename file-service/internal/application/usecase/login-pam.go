package usecase

import (
	"errors"
	"log/slog"

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
	Execute(username, password string) (*jwt.Token, error)
}

type loginPAMUseCase struct {
	config *config.Config
}

var _ LoginPAMUseCase = (*loginPAMUseCase)(nil)

func NewLoginPAMUseCase(config *config.Config) *loginPAMUseCase {
	return &loginPAMUseCase{config: config}
}

func (l *loginPAMUseCase) Execute(username, password string) (*jwt.Token, error) {
	if !l.config.Auth.PAMEnabled {
		return nil, ErrPAMNotEnabled
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
		return nil, ErrPAMFailed
	}

	err = t.Authenticate(0)

	if err != nil {
		slog.Warn("auth attempt failed", "username", username)
		return nil, ErrNotAuthenticated
	}

	// TODO: generate JWT token
	return nil, nil
}
