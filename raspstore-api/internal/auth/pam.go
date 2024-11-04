package auth

import (
	"errors"
	"log/slog"

	"github.com/msteinert/pam"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
)

var (
	ErrUnknowMessageStyle = errors.New("unknown pam message style")
	ErrPAMFailed          = errors.New("pam authentication failed")
	ErrNotAuthenticated   = errors.New("not authenticated")
	ErrPAMNotEnabled      = errors.New("pam not enable for this domain")
)

func LoginPAM(config *config.Config, username, password string) (string, error) {
	if !config.Auth.PAMEnabled {
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

	privateKey, err := readPrivateKey(config.Storage.Path)

	if err != nil {
		slog.Error("could not read Private Key", "err", err)
		return "", ErrNotAuthenticated
	}

	token, err := generateJWT(username, privateKey)

	if err != nil {
		slog.Error("could not generate JWT", "err", nil)
		return "", ErrNotAuthenticated
	}

	return token, nil
}
