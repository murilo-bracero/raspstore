package handler_test

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

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	rMiddleware "github.com/murilo-bracero/raspstore/commons/pkg/middleware"
	"github.com/murilo-bracero/raspstore/file-info-service/proto/v1/file-info-service/pb"
	"github.com/stretchr/testify/assert"
	v1 "raspstore.github.io/fs-service/api/v1"
	"raspstore.github.io/fs-service/internal"
	"raspstore.github.io/fs-service/internal/api/handler"
	"raspstore.github.io/fs-service/internal/model"
)

const testFilename = "test.txt"
const defaultPath = "/folder1"
const defaultUserId = "e9e28c79-a5e8-4545-bd32-e536e690bd4a"

func TestUploadFileSuccess(t *testing.T) {
	if err := godotenv.Load("../../test.env"); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	useCase := &uploadFileUseCaseMock{}
	ctr := handler.NewFileServeHandler(useCase, nil)

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

	err = writer.WriteField("path", defaultPath)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/files", body)
	ctx := context.WithValue(req.Context(), rMiddleware.UserIdKey, defaultUserId)
	ctx = context.WithValue(ctx, chiMiddleware.RequestIDKey, defaultUserId)
	req = req.WithContext(ctx)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Upload)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Status code must be 201 Created")

	var res v1.UploadSuccessResponse
	err = json.Unmarshal(rr.Body.Bytes(), &res)
	assert.NoError(t, err)

	assert.NotEmpty(t, res.FileId)
	assert.Equal(t, res.OwnerId, defaultUserId)
	assert.Equal(t, res.Filename, testFilename)
	assert.Equal(t, res.Path, defaultPath)
	os.Remove(internal.StoragePath() + res.FileId)
}

func TestUploadFileBadRequestWithNoPath(t *testing.T) {
	if err := godotenv.Load("../../test.env"); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	useCase := &uploadFileUseCaseMock{}
	ctr := handler.NewFileServeHandler(useCase, nil)

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

	ctx := context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id")
	ctx = context.WithValue(ctx, rMiddleware.UserIdKey, defaultUserId)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Upload)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Status code must be 400 BadRequest")

	var res v1.ErrorResponse
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

	useCase := &uploadFileUseCaseMock{}
	ctr := handler.NewFileServeHandler(useCase, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err := writer.WriteField("path", defaultPath)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/files", body)
	ctx := context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id")
	ctx = context.WithValue(ctx, rMiddleware.UserIdKey, defaultUserId)
	req = req.WithContext(ctx)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Upload)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code, "Status code must be 422")
}

func TestUploadFileInternalServerError(t *testing.T) {
	if err := godotenv.Load("../../test.env"); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	useCase := &uploadFileUseCaseMock{shouldReturnError: true}
	ctr := handler.NewFileServeHandler(useCase, nil)

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

	err = writer.WriteField("path", defaultPath)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/files", body)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	ctx := context.WithValue(req.Context(), chiMiddleware.RequestIDKey, "test-trace-id")
	ctx = context.WithValue(ctx, rMiddleware.UserIdKey, defaultUserId)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(ctr.Upload)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Status code must be 500 InternalServerError")
}

func TestDownloadFileSuccess(t *testing.T) {
	if err := godotenv.Load("../../test.env"); err != nil {
		log.Println("Could not load .env file. Using system variables instead")
	}

	fileId := uuid.NewString()
	useCase := &downloadFileUseCaseMock{}
	ctr := handler.NewFileServeHandler(nil, useCase)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/files/%s", fileId), nil)
	ctx := context.WithValue(req.Context(), rMiddleware.UserIdKey, defaultUserId)
	req = req.WithContext(ctx)

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

type uploadFileUseCaseMock struct {
	shouldReturnError bool
}

func (u *uploadFileUseCaseMock) Execute(ctx context.Context, req *pb.CreateFileMetadataRequest, src io.Reader) (fileMetadata *pb.FileMetadata, error_ error) {
	if u.shouldReturnError {
		return nil, errors.New("generic error")
	}

	return &pb.FileMetadata{
		FileId:   uuid.NewString(),
		Filename: testFilename,
		Path:     defaultPath,
		OwnerId:  defaultUserId,
	}, nil
}

type downloadFileUseCaseMock struct {
	shouldReturnError bool
}

func (d *downloadFileUseCaseMock) Execute(ctx context.Context, fileId string) (downloadRep *model.FileDownloadRepresentation, error_ error) {
	if d.shouldReturnError {
		return nil, errors.New("generic error")
	}

	tempFile, _ := createTempFile()

	file, _ := os.Open(tempFile)

	return &model.FileDownloadRepresentation{
		Filename: testFilename,
		File:     file,
		FileSize: 1000,
	}, nil
}
