package model

import (
	"github.com/murilo-bracero/raspstore-protofiles/authentication/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Credential struct {
	Id            primitive.ObjectID `bson:"_id"`
	UserId        string             `bson:"user_id"`
	Email         string             `bson:"email"`
	Secret        string             `bson:"secret"`
	Hash          string             `bson:"password"`
	Has2FAEnabled bool               `bson:"is_2FA_enabled"`
}

func ConvertToModel(req *pb.CreateCredentialsRequest) *Credential {
	return &Credential{
		UserId:        req.UserId,
		Email:         req.Email,
		Secret:        "",
		Hash:          req.Hash,
		Has2FAEnabled: false,
	}
}
