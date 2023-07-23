package service_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/grpc/server"
	"github.com/murilo-bracero/raspstore/file-info-service/internal/model"
	"github.com/murilo-bracero/raspstore/file-info-service/proto/v1/file-info-service/pb"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestCreateFileMetadataSuccess(t *testing.T) {
	ctx := context.Background()
	uc := &createUseCaseMock{}
	svc := server.NewFileInfoService(nil, uc)

	filename := uuid.NewString()

	res, err := svc.CreateFileMetadata(ctx, &pb.CreateFileMetadataRequest{
		OwnerId:  uuid.NewString(),
		Filename: filename,
		Size:     3214,
	})

	assert.NoError(t, err)

	assert.NotEmpty(t, res.FileId)
	assert.Equal(t, filename, res.Filename)
	assert.NotEmpty(t, res.OwnerId)
}

func TestCreateFileMetadataFail(t *testing.T) {
	ctx := context.Background()
	uc := &createUseCaseMock{shouldReturnErr: true}
	svc := server.NewFileInfoService(nil, uc)

	filename := uuid.NewString()

	_, err := svc.CreateFileMetadata(ctx, &pb.CreateFileMetadataRequest{
		OwnerId:  uuid.NewString(),
		Filename: filename,
	})

	assert.Error(t, err)
}

func TestFindFileMetadataByIdSuccess(t *testing.T) {
	ctx := context.Background()
	uc := &getFileUseCaseMock{}
	svc := server.NewFileInfoService(uc, nil)

	id := primitive.NewObjectID().Hex()

	res, err := svc.FindFileMetadataById(ctx, &pb.FindFileMetadataByIdRequest{
		FileId: id,
	})

	assert.NoError(t, err)

	assert.Equal(t, id, res.FileId)
	assert.NotEmpty(t, res.Filename)
	assert.NotEmpty(t, res.OwnerId)
}

func TestFindFileMetadataByIdFail(t *testing.T) {
	ctx := context.Background()
	uc := &getFileUseCaseMock{shouldReturnErr: true}
	svc := server.NewFileInfoService(uc, nil)

	id := primitive.NewObjectID().Hex()

	_, err := svc.FindFileMetadataById(ctx, &pb.FindFileMetadataByIdRequest{
		FileId: id,
	})

	assert.Error(t, err)
}

type createUseCaseMock struct {
	shouldReturnErr bool
}

func (f *createUseCaseMock) Execute(file *model.File) (error_ error) {
	if f.shouldReturnErr {
		return mongo.ErrClientDisconnected
	}

	return nil
}

type getFileUseCaseMock struct {
	shouldReturnErr bool
}

func (f *getFileUseCaseMock) Execute(userId string, fileId string) (file *model.File, error_ error) {
	if f.shouldReturnErr {
		return nil, mongo.ErrClientDisconnected
	}

	return &model.File{
		FileId:    fileId,
		Filename:  fileId,
		Size:      int64(rand.Int()),
		UpdatedAt: time.Now(),
		CreatedBy: uuid.NewString(),
		UpdatedBy: uuid.NewString(),
	}, nil
}
