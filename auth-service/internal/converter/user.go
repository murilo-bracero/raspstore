package converter

import (
	v1 "github.com/murilo-bracero/raspstore/auth-service/api/v1"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
)

func ToUserRepresentation(u *model.User) *v1.UserRepresentation {
	return &v1.UserRepresentation{
		UserID:        u.UserId,
		Username:      u.Username,
		IsMfaEnabled:  u.IsMfaEnabled,
		IsMfaVerified: u.IsMfaVerified,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
