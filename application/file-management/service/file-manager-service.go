package service

import (
	"bytes"
	"context"
	"io"
	"log"

	"raspstore.github.io/file-manager/model"
	"raspstore.github.io/file-manager/pb"
	"raspstore.github.io/file-manager/repository"
	"raspstore.github.io/file-manager/system"
	"raspstore.github.io/file-manager/validators"
)

type fileManagerService struct {
	fileRepository repository.FilesRepository
	diskStore      system.DiskStore
	pb.UnimplementedFileManagerServiceServer
}

func NewFileManagerService(diskStore system.DiskStore, fileRepository repository.FilesRepository) pb.FileManagerServiceServer {
	return &fileManagerService{fileRepository: fileRepository, diskStore: diskStore}
}

func (f *fileManagerService) Upload(stream pb.FileManagerService_UploadServer) error {

	req, err := stream.Recv()

	if err != nil {
		log.Print("error while receiving streaming ", err.Error())
		return validators.ErrReceiveStreaming
	}

	if err := validators.ValidateUpload(req); err != nil {
		return err
	}

	file := new(model.File)
	file.FromProtoBuffer(req.GetFiledata())

	data := bytes.Buffer{}
	size := 0

	log.Print(file.Id.Hex(), ": receiving chunks from streaming")

	for {

		req, err := stream.Recv()
		if err == io.EOF {
			log.Print(file.Id.Hex(), ": end file")
			break
		}

		if err != nil {
			log.Print(file.Id.Hex(), ": an error occurs upload file :", err.Error())
			return validators.ErrUploadFile
		}

		chunk := req.GetChunk()
		size += len(chunk)

		if _, err := data.Write(chunk); err != nil {
			log.Print(file.Id.Hex(), ": an error occurs while writing file to buffer: ", err.Error())
			return validators.ErrUploadFile
		}
	}

	if err := f.diskStore.Save(file, data); err != nil {
		log.Print(file.Id.Hex(), ": an error occur while writing file to disk: ", err.Error())
		return validators.ErrUploadFile
	}

	file.Size = uint32(size)

	stream.SendAndClose(file.ToProtoBuffer())

	log.Print(file.Id.Hex(), ": file saved successfully")
	return nil
}

func (f *fileManagerService) Download(req *pb.GetFileRequest, stream pb.FileManagerService_DownloadServer) error {
	return nil
}

func (f *fileManagerService) Delete(ctx context.Context, req *pb.GetFileRequest) (*pb.DeleteFileResponse, error) {
	return nil, nil
}

func (f *fileManagerService) Update(stream pb.FileManagerService_UpdateServer) error {
	return nil
}
