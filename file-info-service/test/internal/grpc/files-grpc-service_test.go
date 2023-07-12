package service_test

import (
	"context"
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore-protofiles/file-info-service/pb"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	service "raspstore.github.io/file-manager/internal/grpc"
	"raspstore.github.io/file-manager/internal/model"
)

func TestCreateFileMetadataSuccess(t *testing.T) {
	ctx := context.Background()
	fr := &filesRepositoryMock{}
	svc := service.NewFileInfoService(fr)

	filename := uuid.NewString()

	res, err := svc.CreateFileMetadata(ctx, &pb.CreateFileMetadataRequest{
		OwnerId:  uuid.NewString(),
		Filename: filename,
		Path:     uuid.NewString() + "/" + filename,
		Size:     3214,
	})

	assert.NoError(t, err)

	assert.NotEmpty(t, res.FileId)
	assert.Equal(t, filename, res.Filename)
	assert.NotEmpty(t, res.OwnerId)
	assert.NotEmpty(t, res.Path)
}

func TestCreateFileMetadataFail(t *testing.T) {
	ctx := context.Background()
	fr := &filesRepositoryMock{shouldReturnErr: true}
	svc := service.NewFileInfoService(fr)

	filename := uuid.NewString()

	_, err := svc.CreateFileMetadata(ctx, &pb.CreateFileMetadataRequest{
		OwnerId:  uuid.NewString(),
		Filename: filename,
		Path:     uuid.NewString() + "/" + filename,
	})

	assert.Error(t, err)
}

func TestFindFileMetadataByIdSuccess(t *testing.T) {
	ctx := context.Background()
	fr := &filesRepositoryMock{}
	svc := service.NewFileInfoService(fr)

	id := primitive.NewObjectID().Hex()

	res, err := svc.FindFileMetadataById(ctx, &pb.FindFileMetadataByIdRequest{
		FileId: id,
	})

	assert.NoError(t, err)

	assert.Equal(t, id, res.FileId)
	assert.NotEmpty(t, res.Filename)
	assert.NotEmpty(t, res.OwnerId)
	assert.NotEmpty(t, res.Path)
}

type filesRepositoryMock struct {
	shouldReturnErr bool
}

func (f *filesRepositoryMock) Save(file *model.File) error {
	if f.shouldReturnErr {
		return mongo.ErrClientDisconnected
	}

	return nil
}

func (f *filesRepositoryMock) FindById(userId string, id string) (*model.File, error) {
	if f.shouldReturnErr {
		return nil, mongo.ErrClientDisconnected
	}

	return &model.File{
		FileId:    id,
		Filename:  id,
		Path:      uuid.NewString() + "/" + id,
		Size:      int64(rand.Int()),
		UpdatedAt: time.Now(),
		CreatedBy: uuid.NewString(),
		UpdatedBy: uuid.NewString(),
	}, nil
}

func (f *filesRepositoryMock) FindByIdLookup(userId string, id string) (fileMetadata *model.FileMetadataLookup, err error) {
	return nil, errors.New("Not Implemented!")
}

func (f *filesRepositoryMock) Delete(userId string, fileId string) error {
	return errors.New("Not Implemented!")
}

func (f *filesRepositoryMock) Update(userId string, file *model.File) error {
	return errors.New("Not Implemented!")
}

func (f *filesRepositoryMock) FindAll(userId string, page int, size int) (filesPage *model.FilePage, err error) {
	return nil, errors.New("Not Implemented!")
}
