package service

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"

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
	file.FromCreateProto(req.GetFiledata())

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

	file.Size = uint32(size)

	if err := f.diskStore.Save(file, data); err != nil {
		log.Print(file.Id.Hex(), ": an error occur while writing file to disk: ", err.Error())
		return validators.ErrUploadFile
	}

	if err := f.fileRepository.Save(file); err != nil {
		log.Println(file.Id, ": an error occur while saving file to database: ", err.Error())
		return err
	}

	stream.SendAndClose(file.ToProtoBuffer())

	log.Print(file.Id.Hex(), ": file saved successfully")
	return nil
}

func (f *fileManagerService) Download(req *pb.GetFileRequest, stream pb.FileManagerService_DownloadServer) error {
	if err := validators.ValidateDownload(req); err != nil {
		return nil
	}

	file, err := f.fileRepository.FindById(req.Id)

	if err != nil {
		return err
	}

	var blob *os.File
	blob, err = os.Open(file.Uri)

	if err != nil {
		return err
	}

	defer blob.Close()

	buff := make([]byte, 64*1024)

	for {
		read, err := blob.Read(buff)

		if err == io.EOF {
			log.Println("end file ", file.Id)
			break
		}

		if err != nil {
			log.Println(file.Id, ": an error occured while serving file: ", err.Error())
			return err
		}

		res := &pb.File{
			Info: &pb.FileMetadata{
				Uri:       file.Uri,
				UpdatedAt: file.UpdatedAt.Unix(),
				CreatedBy: file.CreatedBy,
				UpdatedBy: file.UpdatedBy,
			},
			Chunk: buff[:read],
		}
		if err := stream.Send(res); err != nil {
			log.Panicln(file.Id, ": error sending file to stream: ", err.Error())
			return err
		}

	}

	return nil
}

func (f *fileManagerService) Delete(ctx context.Context, req *pb.GetFileRequest) (*pb.DeleteFileResponse, error) {
	if err := validators.ValidateDownload(req); err != nil {
		return nil, err
	}

	file, err := f.fileRepository.FindById(req.Id)

	if err != nil {
		return nil, err
	}

	if err := f.diskStore.Delete(file.Uri); err != nil {
		return nil, err
	}

	return &pb.DeleteFileResponse{Id: file.Id.Hex()}, nil

}

func (f *fileManagerService) Update(stream pb.FileManagerService_UpdateServer) error {
	req, err := stream.Recv()

	if err != nil {
		log.Print("error while receiving streaming ", err.Error())
		return validators.ErrReceiveStreaming
	}

	if err := validators.ValidateUpdate(req); err != nil {
		return err
	}

	file := new(model.File)
	file.FromUpdateProto(req.GetFiledata())

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

	file.Size = uint32(size)

	if err := f.fileRepository.Update(file); err != nil {
		return err
	}

	stream.SendAndClose(file.ToProtoBuffer())

	log.Print(file.Id.Hex(), ": file saved successfully")
	return nil
}

func (f *fileManagerService) ListFiles(req *pb.FindAllRequest, stream pb.FileManagerService_ListFilesServer) error {
	files, err := f.fileRepository.FindAll()

	if err != nil {
		return err
	}

	for _, file := range files {
		if err := stream.Send(file.ToProtoBuffer()); err != nil {
			return err
		}
	}
	return nil
}
