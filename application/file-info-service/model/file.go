package model

import (
	"time"

	"github.com/murilo-bracero/raspstore-protofiles/file-info-service/pb"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (f *File) ToProtoBuffer() *pb.FileMetadata {
	return &pb.FileMetadata{
		FileId:   f.Id.Hex(),
		Filename: f.Filename,
		Path:     f.Uri,
		OwnerId:  f.CreatedBy,
	}
}

func (f *File) FromCreateProto(file *pb.CreateFileMetadataRequest) {
	f.Id = primitive.NewObjectID()
	f.UpdatedAt = time.Now()
	f.CreatedBy = file.OwnerId
	f.UpdatedBy = file.OwnerId
	f.Filename = file.Filename
	f.Uri = file.Path
}
