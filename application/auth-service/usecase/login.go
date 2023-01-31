package usecase

import (
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"raspstore.github.io/auth-service/repository"
	"raspstore.github.io/auth-service/token"
	"raspstore.github.io/auth-service/utils"
)

type LoginUseCase interface {
	AuthenticateCredentials(username string, rawPassword string, mfaToken string) (accessToken string, refreshToken string, err error)
}

type loginUseCase struct {
	tokenManager     token.TokenManager
	usersRespository repository.UsersRepository
}

const empty = ""

func NewLoginUseCase(tm token.TokenManager, ur repository.UsersRepository) LoginUseCase {
	return &loginUseCase{tokenManager: tm, usersRespository: ur}
}

func (ls *loginUseCase) AuthenticateCredentials(username string, rawPassword string, mfaToken string) (accessToken string, refreshToken string, err error) {
	usr, err := ls.usersRespository.FindByUsername(username)

	if err != nil {
		return empty, empty, err
	}

	if !isValidPassword(rawPassword, usr.Password) {
		return empty, empty, utils.ErrIncorrectCredentials
	}

	if err != nil {
		return empty, empty, utils.ErrUserNotFound
	}

	if usr.IsMfaEnabled && usr.IsMfaVerified && !isValidTotp(mfaToken, usr.Secret) {
		return empty, empty, utils.ErrIncorrectCredentials
	}

	if accessToken, err = ls.tokenManager.Generate(usr.Id.Hex()); err != nil {
		return empty, empty, err
	}

	if refreshToken, err = createRefreshToken(); err != nil {
		return empty, empty, err
	}

	usr.RefreshToken = refreshToken

	if err = ls.usersRespository.Update(usr); err != nil {
		return empty, empty, err
	}

	return accessToken, refreshToken, nil
}

func createRefreshToken() (refreshToken string, err error) {
	seed := uuid.NewString()

	hash, err := bcrypt.GenerateFromPassword([]byte(seed), bcrypt.DefaultCost)

	if err != nil {
		return empty, err
	}

	return string(hash), nil
}

func isValidPassword(rawPass string, hashPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(rawPass)) == nil
}

func isValidTotp(token string, secret string) bool {
	return totp.Validate(token, secret)
}
