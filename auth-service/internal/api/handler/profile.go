package handler

import (
	"context"
	"net/http"

	cm "github.com/go-chi/chi/v5/middleware"
	v1 "github.com/murilo-bracero/raspstore/auth-service/api/v1"
	"github.com/murilo-bracero/raspstore/auth-service/internal"
	"github.com/murilo-bracero/raspstore/auth-service/internal/api/middleware"
	u "github.com/murilo-bracero/raspstore/auth-service/internal/api/utils"
	"github.com/murilo-bracero/raspstore/auth-service/internal/converter"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
	"github.com/murilo-bracero/raspstore/auth-service/internal/usecase"
	"github.com/murilo-bracero/raspstore/commons/pkg/logger"
)

type ProfileHandler interface {
	GetProfile(w http.ResponseWriter, r *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
	DeleteProfile(w http.ResponseWriter, r *http.Request)
}

type profileHandler struct {
	getUserUseCase    usecase.GetUserUseCase
	updateUserUseCase usecase.UpdateProfileUseCase
	deleteUseCase     usecase.DeleteUserUseCase
}

func NewProfileHandler(profileUseCase usecase.GetUserUseCase, updateUserUseCase usecase.UpdateProfileUseCase, deleteUseCase usecase.DeleteUserUseCase) ProfileHandler {
	return &profileHandler{getUserUseCase: profileUseCase, updateUserUseCase: updateUserUseCase, deleteUseCase: deleteUseCase}
}

func (h *profileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserClaimsCtxKey).(*model.UserClaims)

	user, err := h.getUserUseCase.Execute(r.Context(), claims.Subject)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !user.IsEnabled {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	u.Send(w, converter.ToUserRepresentation(user))
}

func (h *profileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	traceId := r.Context().Value(cm.RequestIDKey).(string)
	claims := r.Context().Value(middleware.UserClaimsCtxKey).(*model.UserClaims)

	if res, err := h.isAccountInactive(r.Context(), claims.Subject); res {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var req v1.UpdateProfileRepresentation
	if err := u.ParseBody(r.Body, &req); err != nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}

	if req.Username == "" {
		u.HandleBadRequest(w, "ERR001", "username could not be null or empty", traceId)
		return
	}

	user := &model.User{
		UserId:   claims.Subject,
		Username: req.Username,
	}

	if err := h.updateUserUseCase.Execute(r.Context(), user); err != nil {
		handleUpdateUserError(w, err)
		return
	}

	u.Send(w, converter.ToUserRepresentation(user))
}

func (h *profileHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value(middleware.UserClaimsCtxKey).(*model.UserClaims)

	if res, err := h.isAccountInactive(r.Context(), claims.Subject); res {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := h.deleteUseCase.Execute(r.Context(), claims.Subject); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	u.NoContent(w)
}

func (h *profileHandler) isAccountInactive(ctx context.Context, userId string) (bool, error) {
	traceId := ctx.Value(cm.RequestIDKey).(string)
	usr, err := h.getUserUseCase.Execute(ctx, userId)

	if err != nil {
		logger.Error("[%s] Could not retrieve user information from the database: userId=%s: %s", traceId, userId, err.Error())
		return false, nil
	}

	return !usr.IsEnabled, nil
}

func handleUpdateUserError(w http.ResponseWriter, err error) {
	if err == internal.ErrUserNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	if err == internal.ErrConflict {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
