package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore-protofiles/file-info-service/pb"
	"github.com/stretchr/testify/assert"
	"raspstore.github.io/fs-service/api/controller"
	"raspstore.github.io/fs-service/api/dto"
	"raspstore.github.io/fs-service/internal"
)

const testFilename = "test.txt"
const defaultUserId = "e9e28c79-a5e8-4545-bd32-e536e690bd4a"

func TestUploadFileSuccess(t *testing.T) {
	if err := godotenv.Load("../../test.env"); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	path := "/folder1"
	useCase := &fileInfoUseCaseMock{defaultOwnerId: defaultUserId}
	ctr := controller.NewFileServeController(useCase)

	tempFile, err := createTempFile()
	assert.NoError(t, err)

	file, err := os.Open(tempFile)
	assert.NoError(t, err)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(tempFile))
	assert.NoError(t, err)

	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	err = writer.WriteField("path", path)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/files", body)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Upload)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Status code must be 201 Created")

	var res dto.UploadSuccessResponse
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	assert.NoError(t, err)

	assert.NotEmpty(t, res.FileId)
	assert.Equal(t, res.OwnerId, defaultUserId)
	assert.Equal(t, res.Filename, testFilename)
	assert.Equal(t, res.Path, path)
	os.Remove(internal.StoragePath() + res.FileId)
}

func TestUploadFileBadRequestWithNoPath(t *testing.T) {
	if err := godotenv.Load("../../test.env"); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	useCase := &fileInfoUseCaseMock{defaultOwnerId: defaultUserId}
	ctr := controller.NewFileServeController(useCase)

	tempFile, err := createTempFile()
	assert.NoError(t, err)

	file, err := os.Open(tempFile)
	assert.NoError(t, err)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(tempFile))
	assert.NoError(t, err)

	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/files", body)
	assert.NoError(t, err)

	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Upload)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code must be 400 BadRequest")

	var res dto.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	assert.NoError(t, err)

	assert.NotEmpty(t, res.Code)
	assert.NotEmpty(t, res.Message)
	assert.NotEmpty(t, res.TraceId)
}

func TestUploadFileBadRequestWithNoFile(t *testing.T) {
	if err := godotenv.Load("../../test.env"); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	path := "/folder1"
	useCase := &fileInfoUseCaseMock{defaultOwnerId: defaultUserId}
	ctr := controller.NewFileServeController(useCase)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err := writer.WriteField("path", path)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/files", body)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Upload)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code must be 400 BadRequest")

	var res dto.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	assert.NoError(t, err)

	assert.NotEmpty(t, res.Code)
	assert.NotEmpty(t, res.Message)
	assert.NotEmpty(t, res.TraceId)
}

func TestUploadFileInternalServerError(t *testing.T) {
	if err := godotenv.Load("../../test.env"); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	path := "/folder1"
	useCase := &fileInfoUseCaseMock{defaultOwnerId: defaultUserId, shouldReturnError: true}
	ctr := controller.NewFileServeController(useCase)

	tempFile, err := createTempFile()
	assert.NoError(t, err)

	file, err := os.Open(tempFile)
	assert.NoError(t, err)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(tempFile))
	assert.NoError(t, err)

	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	err = writer.WriteField("path", path)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/files", body)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-trace-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Upload)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code must be 500 InternalServerError")

	var res dto.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	assert.NoError(t, err)

	assert.NotEmpty(t, res.Code)
	assert.NotEmpty(t, res.Message)
	assert.NotEmpty(t, res.TraceId)
}

func TestDownloadFileSuccess(t *testing.T) {
	if err := godotenv.Load("../../test.env"); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	fileId := uuid.NewString()
	useCase := &fileInfoUseCaseMock{defaultOwnerId: defaultUserId, defaultFileId: fileId}
	ctr := controller.NewFileServeController(useCase)

	f, err := os.Create(internal.StoragePath() + "/" + fileId)
	defer os.Remove(internal.StoragePath() + "/" + fileId)
	assert.NoError(t, err)

	_, err = f.WriteString("Test file")
	assert.NoError(t, err)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/files/%s", fileId), nil)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Download)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200 OK")
	assert.Equal(t, "application/octet-stream", rr.Header().Get("Content-Type"))
	assert.Equal(t, fmt.Sprintf("attachment; filename=\"%s\"", testFilename), rr.Header().Get("Content-Disposition"))
}

func createTempFile() (string, error) {
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, testFilename)
	err := ioutil.WriteFile(tempFile, []byte("test content"), 0666)
	return tempFile, err
}

type fileInfoUseCaseMock struct {
	defaultOwnerId    string
	defaultFileId     string
	shouldReturnError bool
}

func (f *fileInfoUseCaseMock) GetFileMetadataById(id string) (fileMtadata *pb.FileMetadata, err error) {
	if f.shouldReturnError {
		return nil, errors.New("generic error")
	}

	return &pb.FileMetadata{
		FileId:   f.defaultFileId,
		Filename: testFilename,
		Path:     "/",
		OwnerId:  f.defaultOwnerId,
	}, nil
}

func (f *fileInfoUseCaseMock) CreateFileMetadata(req *pb.CreateFileMetadataRequest) (fileMtadata *pb.FileMetadata, err error) {
	if f.shouldReturnError {
		return nil, errors.New("generic error")
	}

	return &pb.FileMetadata{
		FileId:   uuid.NewString(),
		Filename: req.Filename,
		Path:     req.Path,
		OwnerId:  req.OwnerId,
	}, nil
}
