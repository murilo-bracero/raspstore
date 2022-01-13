package model

type Credential struct {
	Id            string `bson:"user_id"`
	Email         string `bson:"email"`
	Secret        string `bson:"secret"`
	Hash          string `bson:"password"`
	Has2FAEnabled bool   `bson:"has_2FA_enabled"`
}
