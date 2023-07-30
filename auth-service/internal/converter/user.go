package converter

import (
	v1 "github.com/murilo-bracero/raspstore/auth-service/api/v1"
	"github.com/murilo-bracero/raspstore/auth-service/internal/model"
)

func ToUserRepresentation(u *model.User) *v1.UserRepresentation {
	return &v1.UserRepresentation{
		UserID:        u.UserId,
		Username:      u.Username,
		IsEnabled:     u.IsEnabled,
		IsMfaEnabled:  u.IsMfaEnabled,
		IsMfaVerified: u.IsMfaVerified,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func ToUser(req *v1.CreateUserRepresentation) *model.User {
	return &model.User{
		Username:     req.Username,
		Password:     req.Password,
		Permissions:  req.Roles,
		IsEnabled:    true,
		IsMfaEnabled: false,
	}
}

func ToPageRepresentation(up *model.UserPage, page int, size int, nextUrl string) v1.PageRepresentation {
	content := make([]*v1.UserRepresentation, len(up.Content))
	for i, usr := range up.Content {
		content[i] = ToUserRepresentation(usr)
	}

	return v1.PageRepresentation{
		Page:          page,
		Size:          size,
		TotalElements: up.Count,
		Next:          nextUrl,
		Content:       content,
	}
}
