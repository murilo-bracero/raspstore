package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"raspstore.github.io/file-manager/pb"
)

type File struct {
	Id        primitive.ObjectID `bson:"_id"`
	Filename  string             `bson:"filename"`
	Uri       string             `bson:"uri"`
	Size      uint32             `bson:"size"`
	UpdatedAt time.Time          `bson:"updated_at"`
	CreatedBy string             `bson:"created_by"`
	UpdatedBy string             `bson:"updated_by"`
}

func (f *File) ToProtoBuffer() *pb.FileRef {
	return &pb.FileRef{
		Id:        f.Id.Hex(),
		Uri:       f.Uri,
		Size:      f.Size,
		UpdatedAt: f.UpdatedAt.Unix(),
		CreatedBy: f.CreatedBy,
		UpdatedBy: f.UpdatedBy,
	}
}

func (f *File) FromProtoBuffer(user *pb.CreateFileRequestData) {
	f.Id = primitive.NewObjectID()
	f.UpdatedAt = time.Now()
	f.CreatedBy = user.CreatedBy
	f.UpdatedBy = user.CreatedBy
	f.Filename = user.Filename
}
