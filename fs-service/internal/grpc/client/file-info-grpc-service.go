package client

import (
	"context"

	"github.com/murilo-bracero/raspstore/file-info-service/proto/v1/file-info-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"raspstore.github.io/fs-service/internal"
)

type FileInfoService interface {
	GetFileMetadataById(fileId string, userId string) (fileMtadata *pb.FileMetadata, err error)
	CreateFileMetadata(req *pb.CreateFileMetadataRequest) (fileMtadata *pb.FileMetadata, err error)
}

type fileInfoService struct {
	ctx context.Context
}

func NewFileInfoService(ctx context.Context) FileInfoService {
	return &fileInfoService{ctx: ctx}
}

func (f *fileInfoService) GetFileMetadataById(fileId string, userId string) (fileMtadata *pb.FileMetadata, err error) {
	conn, err := grpc.Dial(internal.FileInfoServiceUrl(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := pb.NewFileInfoServiceClient(conn)

	return client.FindFileMetadataById(f.ctx, &pb.FindFileMetadataByIdRequest{FileId: fileId, RequesterUserId: userId})
}

func (f *fileInfoService) CreateFileMetadata(req *pb.CreateFileMetadataRequest) (fileMtadata *pb.FileMetadata, err error) {
	conn, err := grpc.Dial(internal.FileInfoServiceUrl(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := pb.NewFileInfoServiceClient(conn)

	return client.CreateFileMetadata(f.ctx, req)
}
