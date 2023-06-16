package service

import (
	"log"

	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"raspstore.github.io/auth-service/internal"
	"raspstore.github.io/auth-service/internal/model"
	"raspstore.github.io/auth-service/internal/repository"
)

type LoginService interface {
	AuthenticateCredentials(username string, rawPassword string, mfaToken string) (accessToken string, refreshToken string, err error)
}

type loginService struct {
	tokenService     TokenService
	usersRespository repository.UsersRepository
}

func NewLoginService(ts TokenService, ur repository.UsersRepository) LoginService {
	return &loginService{tokenService: ts, usersRespository: ur}
}

func (ls *loginService) AuthenticateCredentials(username string, rawPassword string, mfaToken string) (accessToken string, refreshToken string, err error) {
	usr, err := ls.usersRespository.FindByUsername(username)

	if err != nil {
		log.Printf("[ERROR] Could not find user: %s in database due to error: %s", username, err.Error())
		return "", "", err
	}

	if !isValidPassword(rawPassword, usr.Password) {
		return "", "", internal.ErrIncorrectCredentials
	}

	if err := isValidMfa(usr, mfaToken); err != nil {
		return "", "", err
	}

	if accessToken, err = ls.tokenService.Generate(usr); err != nil {
		return "", "", err
	}

	if refreshToken, err = createRefreshToken(); err != nil {
		return "", "", err
	}

	usr.RefreshToken = refreshToken

	if err = ls.usersRespository.Update(usr); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func isValidMfa(usr *model.User, mfaToken string) error {
	if usr.IsMfaEnabled && usr.IsMfaVerified && !isValidTotp(mfaToken, usr.Secret) {
		return internal.ErrIncorrectCredentials
	}

	return nil
}

func createRefreshToken() (refreshToken string, err error) {
	seed := uuid.NewString()

	hash, err := bcrypt.GenerateFromPassword([]byte(seed), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func isValidPassword(rawPass string, hashPass string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(rawPass)) == nil
}

func isValidTotp(token string, secret string) bool {
	return totp.Validate(token, secret)
}
