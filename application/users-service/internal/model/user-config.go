package model

type UserConfiguration struct {
	StorageLimit            string `bson:"storage_limit"`
	MinPasswordLength       int    `bson:"min_password_length"`
	AllowPublicUserCreation bool   `bson:"allow_public_user_creation"`
	AllowLoginWithEmail     bool   `bson:"allow_login_with_email"`
	EnforceMfa              bool   `bson:"enforce_mfa"`
}
