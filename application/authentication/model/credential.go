package model

type Credential struct {
	Id    string `bson:"user_id"`
	Email string `bson:"email"`
	Hash  string `bson:"password"`
}
