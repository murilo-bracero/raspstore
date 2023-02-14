package usecase

import (
	"context"

	"github.com/murilo-bracero/raspstore-protofiles/file-info-service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"raspstore.github.io/fs-service/internal"
)

type FileInfoUseCase interface {
	GetFileMetadataById(id string) (fileMtadata *pb.FileMetadata, err error)
	CreateFileMetadata(req *pb.CreateFileMetadataRequest) (fileMtadata *pb.FileMetadata, err error)
}

type fileInfoUseCase struct {
	ctx context.Context
}

func NewFileInfoUseCase(ctx context.Context) FileInfoUseCase {
	return &fileInfoUseCase{ctx: ctx}
}

func (f *fileInfoUseCase) GetFileMetadataById(id string) (fileMtadata *pb.FileMetadata, err error) {
	conn, err := grpc.Dial(internal.FileInfoServiceUrl(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := pb.NewFileInfoServiceClient(conn)

	return client.FindFileMetadataById(f.ctx, &pb.FindFileMetadataByIdRequest{FileId: id})
}

func (f *fileInfoUseCase) CreateFileMetadata(req *pb.CreateFileMetadataRequest) (fileMtadata *pb.FileMetadata, err error) {
	conn, err := grpc.Dial(internal.FileInfoServiceUrl(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	client := pb.NewFileInfoServiceClient(conn)

	return client.CreateFileMetadata(f.ctx, req)
}
