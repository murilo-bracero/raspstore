package model

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"raspstore.github.io/file-manager/pb"
)

type File struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Filename  string             `json:"filename,omitempty" bson:"filename"`
	Uri       string             `json:"-" bson:"uri"`
	Size      uint32             `json:"size,omitempty" bson:"size"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at"`
	CreatedBy string             `json:"created_by,omitempty" bson:"created_by"`
	UpdatedBy string             `json:"updated_by,omitempty" bson:"updated_by"`
}

func NewFile(filename string, createdBy string, size uint32) *File {
	return &File{
		Id:        primitive.NewObjectID(),
		Filename:  filename,
		CreatedBy: createdBy,
		UpdatedBy: createdBy,
		Size:      size,
		UpdatedAt: time.Now(),
	}
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
