package usecase

import (
	"time"

	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/idp/internal"
	"github.com/murilo-bracero/raspstore/idp/internal/infra"
	"github.com/murilo-bracero/raspstore/idp/internal/model"
	"github.com/murilo-bracero/raspstore/idp/internal/repository"
	"github.com/murilo-bracero/raspstore/idp/internal/token"
)

type LoginUseCase interface {
	AuthenticateCredentials(username string, rawPassword string, mfaToken string) (tokenCredentials *model.TokenCredentials, err error)
}

type loginUseCase struct {
	usersRespository repository.UsersRepository
	config           *infra.Config
}

func NewLoginUseCase(config *infra.Config, ur repository.UsersRepository) LoginUseCase {
	return &loginUseCase{usersRespository: ur, config: config}
}

func (ls *loginUseCase) AuthenticateCredentials(username string, rawPassword string, mfaToken string) (tokenCredentials *model.TokenCredentials, err error) {
	usr, err := ls.usersRespository.FindByUsername(username)

	if err != nil {
		logger.Error("Could not find user: %s in database due to error: %s", username, err.Error())
		return nil, err
	}

	if !usr.IsEnabled {
		return nil, internal.ErrInactiveAccount
	}

	if !usr.Authenticate(rawPassword) {
		return nil, internal.ErrIncorrectCredentials
	}

	if err := usr.ValidateTotpToken(mfaToken); err != nil {
		return nil, err
	}

	tokenCredentials = &model.TokenCredentials{}

	if tokenCredentials.AccessToken, err = token.Generate(ls.config, usr); err != nil {
		return nil, err
	}

	if err = usr.GenerateRefreshToken(); err != nil {
		return nil, err
	}

	tokenCredentials.RefreshToken = usr.RefreshToken

	if err = ls.usersRespository.Update(usr); err != nil {
		return nil, err
	}

	tokenCredentials.ExpirestAt = time.Now().Add(time.Duration(ls.config.TokenDuration) * time.Second)

	return tokenCredentials, nil
}
