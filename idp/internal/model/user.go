package model

import (
	"time"

	"github.com/google/uuid"
	v1 "github.com/murilo-bracero/raspstore/idp/api/v1"
	"github.com/murilo-bracero/raspstore/idp/internal"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UserId        string `bson:"user_id"`
	Username      string
	IsEnabled     bool   `bson:"is_enabled"`
	PasswordHash  string `bson:"password_hash"`
	Secret        string
	Permissions   []string
	RefreshToken  string    `bson:"refresh_token"`
	IsMfaEnabled  bool      `bson:"is_mfa_enabled"`
	IsMfaVerified bool      `bson:"is_mfa_verified"`
	CreatedAt     time.Time `bson:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at"`
}

func (u *User) ToUserRepresentation() *v1.UserRepresentation {
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

func NewUser(req *v1.CreateUserRepresentation) *User {
	user := &User{
		Username:     req.Username,
		Permissions:  req.Roles,
		IsEnabled:    true,
		IsMfaEnabled: false,
	}

	user.SetPassword(req.Password)
	return user
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.PasswordHash = string(hash)
	return nil
}

func (u *User) Authenticate(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) == nil
}

func (u *User) GenerateRefreshToken() error {
	seed := uuid.NewString()

	hash, err := bcrypt.GenerateFromPassword([]byte(seed), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.RefreshToken = string(hash)
	return nil
}

func (u *User) ValidateTotpToken(token string) error {
	if u.IsMfaEnabled && u.IsMfaVerified && !totp.Validate(token, u.Secret) {
		return internal.ErrIncorrectCredentials
	}

	return nil
}

func (u *User) BeforeCreate() {
	u.UserId = uuid.NewString()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}
