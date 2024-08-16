package main_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/entity"
	"github.com/murilo-bracero/raspstore/file-service/internal/domain/model"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ApiTest struct {
	KeycloakUrl string
	ApiUrl      string
	keycloakC   tc.Container
	apiC        tc.Container
}

func NewApiTest(ctx context.Context) (*ApiTest, error) {
	net, err := network.New(ctx)

	if err != nil {
		return nil, err
	}

	kcr := tc.ContainerRequest{
		Image:        "quay.io/keycloak/keycloak:25.0",
		ExposedPorts: []string{"8080/tcp"},
		WaitingFor:   wait.ForLog(`Running the server.*development mode\.`).AsRegexp(),
		Env: map[string]string{
			"KEYCLOAK_ADMIN":          "admin",
			"KEYCLOAK_ADMIN_PASSWORD": "admin_test_password001",
		},
		Cmd:      []string{"start-dev"},
		Networks: []string{net.Name},
		NetworkAliases: map[string][]string{
			net.Name: {"keycloak.rstore.com"},
		},
	}
	kc, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: kcr,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	keycloakUrl, err := getContainerUrl(ctx, kc, "8080")

	if err != nil {
		return nil, err
	}

	acr := tc.ContainerRequest{
		FromDockerfile: tc.FromDockerfile{
			Context:    "..",
			Dockerfile: "Dockerfile",
		},
		ExposedPorts: []string{"9090/tcp"},
		Env: map[string]string{
			"PUBLIC_KEY_URL": "http://keycloak.rstore.com:8080/realms/master/protocol/openid-connect/certs",
		},
		Networks: []string{net.Name},
	}
	ac, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: acr,
		Started:          true,
	})

	if err != nil {
		return nil, err
	}

	apiUrl, err := getContainerUrl(ctx, ac, "9090")

	if err != nil {
		return nil, err
	}

	return &ApiTest{KeycloakUrl: keycloakUrl, ApiUrl: apiUrl, keycloakC: kc, apiC: ac}, nil
}

func (a *ApiTest) Cleanup(ctx context.Context) {
	a.apiC.Terminate(ctx)
	a.keycloakC.Terminate(ctx)
}

func getContainerUrl(ctx context.Context, c tc.Container, port nat.Port) (string, error) {
	hst, err := c.Host(ctx)

	if err != nil {
		return "", err
	}

	prt, err := c.MappedPort(ctx, port)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s:%d", hst, prt.Int()), nil
}

func TestService(t *testing.T) {
	ctx := context.Background()

	apiTest, err := NewApiTest(ctx)

	assert.NoError(t, err, "NewApiTest")

	t.Cleanup(func() {
		apiTest.Cleanup(ctx)
	})

	token, err := getToken(apiTest.KeycloakUrl)

	assert.NoError(t, err, "getToken")

	t.Run("POST /upload - Upload file successfully should return status code CREATED", func(t *testing.T) {
		resource := apiTest.ApiUrl + "/file-service/v1/uploads"

		tempFile, err := createTempFile("test.txt")
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

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, resource, body)

		assert.NoError(t, err, "NewRequest")

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		res, err := client.Do(req)

		assert.NoError(t, err, "client.Do")

		defer res.Body.Close()

		var response map[string]interface{}
		err = json.NewDecoder(res.Body).Decode(&response)
		assert.NoError(t, err, "NewDecoder")

		assert.Equal(t, http.StatusCreated, res.StatusCode)
		assert.NotEmpty(t, response["fileId"])
		assert.Equal(t, response["filename"], "test.txt")
		assert.NotEmpty(t, response["ownerId"])
	})

	t.Run("POST /upload - Upload fail when file is not provided should return UNPROCESSABLE ENTITY", func(t *testing.T) {
		resource := apiTest.ApiUrl + "/file-service/v1/uploads"

		tempFile, err := createTempFile("test.txt")
		assert.NoError(t, err)

		file, err := os.Open(tempFile)
		assert.NoError(t, err)
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		err = writer.Close()
		assert.NoError(t, err)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPost, resource, body)

		assert.NoError(t, err, "NewRequest")

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		res, err := client.Do(req)

		assert.NoError(t, err, "client.Do")

		defer res.Body.Close()

		assert.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
		assert.NotEmpty(t, res.Header.Get("X-Trace-Id"))
	})

	t.Run("GET /files - Get all files from user should return OK", func(t *testing.T) {
		resource := apiTest.ApiUrl + "/file-service/v1/files"

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, resource, http.NoBody)

		assert.NoError(t, err, "NewRequest")

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")

		res, err := client.Do(req)

		assert.NoError(t, err, "client.Do")

		defer res.Body.Close()

		var fpr model.FilePageResponse
		err = json.NewDecoder(res.Body).Decode(&fpr)
		assert.NoError(t, err, "NewDecoder")

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, 0, fpr.Size)
		assert.Equal(t, 1, fpr.TotalElements)
		assert.Equal(t, 0, fpr.Page)
		assert.Empty(t, fpr.Next)
		assert.NotEmpty(t, fpr.Content)

		fileRef := fpr.Content[0]

		assert.Equal(t, "test.txt", fileRef.Filename)
		assert.NotEmpty(t, fileRef.FileId)
		assert.NotEmpty(t, fileRef.Owner)
		assert.NotEmpty(t, fileRef.Size)
		assert.NotEmpty(t, fileRef.CreatedAt)

		userId := fileRef.Owner

		assert.Equal(t, fileRef.CreatedBy, userId)

		assert.NotEmpty(t, res.Header.Get("X-Trace-Id"))
	})

	t.Run("GET /files - Get all files from user with pagination should return OK", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			_, err := uploadFile(apiTest, token, uuid.NewString())
			assert.NoError(t, err, "uploadFile")
		}

		resource := apiTest.ApiUrl + "/file-service/v1/files?size=10&page=0"

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, resource, http.NoBody)

		assert.NoError(t, err, "NewRequest")

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")

		res, err := client.Do(req)

		assert.NoError(t, err, "client.Do")

		defer res.Body.Close()

		var fpr model.FilePageResponse
		err = json.NewDecoder(res.Body).Decode(&fpr)
		assert.NoError(t, err, "NewDecoder")

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, 10, fpr.Size)
		assert.GreaterOrEqual(t, fpr.TotalElements, 20)
		assert.Equal(t, 0, fpr.Page)
		assert.Equal(t, apiTest.ApiUrl+"/file-service/v1/files?page=1&size=10", "http://"+fpr.Next)
		assert.NotEmpty(t, fpr.Content)
		assert.Equal(t, 10, len(fpr.Content))
	})

	t.Run("GET /files - Get all files from user with pagination and filename query should return OK", func(t *testing.T) {

		_, err := uploadFile(apiTest, token, "queryable_file.txt")
		assert.NoError(t, err, "uploadFile")

		resource := apiTest.ApiUrl + "/file-service/v1/files?size=10&page=0&filename=queryable_file"

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, resource, http.NoBody)

		assert.NoError(t, err, "NewRequest")

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")

		res, err := client.Do(req)

		assert.NoError(t, err, "client.Do")

		defer res.Body.Close()

		var fpr model.FilePageResponse
		err = json.NewDecoder(res.Body).Decode(&fpr)
		assert.NoError(t, err, "NewDecoder")

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, 10, fpr.Size)
		assert.GreaterOrEqual(t, fpr.TotalElements, 1)
		assert.Equal(t, 0, fpr.Page)
		assert.Empty(t, fpr.Next)
		assert.NotEmpty(t, fpr.Content)
		assert.Equal(t, 1, len(fpr.Content))

		fileRef := fpr.Content[0]

		assert.Equal(t, "queryable_file.txt", fileRef.Filename)
	})

	t.Run("PUT /files - Update file by ID should return OK", func(t *testing.T) {
		fc, err := uploadFile(apiTest, token, "queryable_file.txt")
		assert.NoError(t, err, "uploadFile")

		resource := fmt.Sprintf("%s/file-service/v1/files/%s", apiTest.ApiUrl, fc.FileId)

		body := strings.NewReader(`
			{
				"filename" : "updated_file_name.txt",
				"secret"   : true
			}`)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, resource, body)

		assert.NoError(t, err, "NewRequest")

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)

		assert.NoError(t, err, "client.Do")

		defer res.Body.Close()

		var f entity.File
		err = json.NewDecoder(res.Body).Decode(&f)
		assert.NoError(t, err, "NewDecoder")

		assert.Equal(t, "updated_file_name.txt", f.Filename)
		assert.NotEmpty(t, f.UpdatedAt)
		assert.NotEmpty(t, f.UpdatedBy)
		assert.NotEqual(t, f.CreatedAt, f.UpdatedAt)
		assert.Equal(t, "updated_file_name.txt", f.Filename)
	})

	t.Run("PUT /files - Update file with invalid payload should return BAD REQUEST", func(t *testing.T) {

		fc, err := uploadFile(apiTest, token, uuid.NewString())
		assert.NoError(t, err, "pickOneFile")

		resource := fmt.Sprintf("%s/file-service/v1/files/%s", apiTest.ApiUrl, fc.FileId)

		body := strings.NewReader(`
			{
				"unrecognizable" : "updated_file_name.txt"
			}`)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, resource, body)

		assert.NoError(t, err, "NewRequest")

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)

		assert.NoError(t, err, "client.Do")

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	})

	t.Run("DELETE /files - Delete file by ID should return NO CONTENT", func(t *testing.T) {

		fc, err := uploadFile(apiTest, token, uuid.NewString())
		assert.NoError(t, err, "pickOneFile")

		resource := fmt.Sprintf("%s/file-service/v1/files/%s", apiTest.ApiUrl, fc.FileId)

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodDelete, resource, http.NoBody)

		assert.NoError(t, err, "NewRequest")

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		res, err := client.Do(req)

		assert.NoError(t, err, "client.Do")

		assert.Equal(t, http.StatusNoContent, res.StatusCode)

		_, err = findFileById(apiTest, token, fc.FileId)
		assert.Error(t, err, "findFileById")

		assert.Equal(t, "404 Not Found", err.Error())
	})
}

func getToken(keycloakUrl string) (string, error) {

	tokenUrl := keycloakUrl + "/realms/master/protocol/openid-connect/token"

	payload := strings.NewReader("username=admin&password=admin_test_password001&grant_type=password&client_id=admin-cli")

	r, err := http.NewRequest(http.MethodPost, tokenUrl, payload)

	if err != nil {
		return "", err
	}

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(r)

	if err != nil {
		slog.Error("Error during request", "err", err)
		return "", err
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("non OK response from kaycloak")
	}

	var response map[string]interface{}
	if err := json.Unmarshal(responseBody, &response); err != nil {
		slog.Error("Error decoding response", "err", err)
		return "", err
	}

	return response["access_token"].(string), nil
}

func createTempFile(filename string) (string, error) {
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, filename)
	err := os.WriteFile(tempFile, []byte("test content"), 0666)
	return tempFile, err
}

func uploadFile(apiTest *ApiTest, token string, filename string) (*model.UploadSuccessResponse, error) {
	resource := apiTest.ApiUrl + "/file-service/v1/uploads"

	tempFile, err := createTempFile(filename)

	if err != nil {
		return nil, err
	}

	file, err := os.Open(tempFile)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(tempFile))

	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, file)

	if err != nil {
		return nil, err
	}

	err = writer.Close()

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, resource, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return nil, errors.New("non ok status")
	}

	var response *model.UploadSuccessResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func findFileById(apiTest *ApiTest, token string, fileId string) (*entity.File, error) {
	resource := fmt.Sprintf("%s/file-service/v1/files/%s", apiTest.ApiUrl, fileId)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, resource, http.NoBody)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	defer res.Body.Close()

	var f *entity.File
	err = json.NewDecoder(res.Body).Decode(&f)

	if err != nil {
		return nil, err
	}

	return f, nil
}
