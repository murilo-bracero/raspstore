package model

import (
	"log"
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

func (f *File) FromCreateProto(file *pb.CreateFileRequestData) {
	f.Id = primitive.NewObjectID()
	f.UpdatedAt = time.Now()
	f.CreatedBy = file.CreatedBy
	f.UpdatedBy = file.CreatedBy
	f.Filename = file.Filename
}

func (f *File) FromUpdateProto(file *pb.UpdateFileRequestData) {
	f.UpdatedAt = time.Now()
	f.UpdatedBy = file.UpdatedBy
	if id, err := primitive.ObjectIDFromHex(file.Id); err != nil {
		log.Println("error converting ", file.Id, " to ObjectId")
	} else {
		f.Id = id
	}

}
