package model

import v1 "github.com/murilo-bracero/raspstore/users-service/api/v1"

type UserConfiguration struct {
	StorageLimit            string `bson:"storage_limit"`
	MinPasswordLength       int    `bson:"min_password_length"`
	AllowPublicUserCreation bool   `bson:"allow_public_user_creation"`
	AllowLoginWithEmail     bool   `bson:"allow_login_with_email"`
	EnforceMfa              bool   `bson:"enforce_mfa"`
}

func (m *UserConfiguration) ToUserConfigurationResponse() v1.UserConfigurationResponse {
	return v1.UserConfigurationResponse{
		StorageLimit:            m.StorageLimit,
		MinPasswordLength:       m.MinPasswordLength,
		AllowPublicUserCreation: m.AllowPublicUserCreation,
		AllowLoginWithEmail:     m.AllowLoginWithEmail,
		EnforceMfa:              m.EnforceMfa,
	}
}
