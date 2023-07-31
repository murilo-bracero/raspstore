package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
	"github.com/murilo-bracero/raspstore/idp/internal"
	"github.com/murilo-bracero/raspstore/idp/internal/infra"
	"github.com/murilo-bracero/raspstore/idp/internal/model"
	"github.com/murilo-bracero/raspstore/idp/internal/repository"
	"github.com/murilo-bracero/raspstore/idp/internal/token"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
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

	if !isValidPassword(rawPassword, usr.Password) {
		return nil, internal.ErrIncorrectCredentials
	}

	if err := isValidMfa(usr, mfaToken); err != nil {
		return nil, err
	}

	tokenCredentials = &model.TokenCredentials{}

	if tokenCredentials.AccessToken, err = token.Generate(ls.config, usr); err != nil {
		return nil, err
	}

	if tokenCredentials.RefreshToken, err = generateRefreshToken(); err != nil {
		return nil, err
	}

	usr.RefreshToken = tokenCredentials.RefreshToken

	if err = ls.usersRespository.Update(usr); err != nil {
		return nil, err
	}

	tokenCredentials.ExpirestAt = time.Now().Add(time.Duration(ls.config.TokenDuration) * time.Second)

	return tokenCredentials, nil
}

func isValidPassword(rawPass string, hashPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(rawPass)) == nil
}

func isValidMfa(usr *model.User, mfaToken string) error {
	if usr.IsMfaEnabled && usr.IsMfaVerified && !totp.Validate(mfaToken, usr.Secret) {
		return internal.ErrIncorrectCredentials
	}

	return nil
}

func generateRefreshToken() (refreshToken string, err error) {
	seed := uuid.NewString()

	hash, err := bcrypt.GenerateFromPassword([]byte(seed), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}
