package handler

import (
	"net/http"
	"strings"

	"github.com/murilo-bracero/raspstore/auth-service/internal"
	u "github.com/murilo-bracero/raspstore/auth-service/internal/api/utils"
	"github.com/murilo-bracero/raspstore/auth-service/internal/converter"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/token"
	"github.com/murilo-bracero/raspstore/auth-service/internal/usecase"
)

type ProfileHandler interface {
	GetProfile(w http.ResponseWriter, r *http.Request)
}

type profileHandler struct {
	profileUseCase usecase.GetProfileUseCase
}

func NewProfileHandler(profileUseCase usecase.GetProfileUseCase) ProfileHandler {
	return &profileHandler{profileUseCase: profileUseCase}
}

func (h *profileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, err := getClaimsForRequest(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	user, err := h.profileUseCase.Execute(r.Context(), claims.Subject)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u.Send(w, converter.ToUserRepresentation(user))
}

func getClaimsForRequest(r *http.Request) (*model.UserClaims, error) {
	claims, err := checkTokenInCookie(r)

	if err == nil {
		return claims, nil
	}

	return checkTokenInHeader(r)
}

func checkTokenInCookie(r *http.Request) (*model.UserClaims, error) {
	accessCookie, err := r.Cookie("access_token")

	if err != nil {
		return nil, err
	}

	accessToken := strings.ReplaceAll(accessCookie.Value, "Bearer ", "")

	return token.Verify(accessToken)
}

func checkTokenInHeader(r *http.Request) (*model.UserClaims, error) {
	header := r.Header.Get("Authorization")

	if header == "" {
		return nil, internal.ErrEmptyToken
	}

	accessToken := strings.ReplaceAll(header, "Bearer ", "")

	return token.Verify(accessToken)
}
