package service_test

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore-protofiles/file-manager/pb"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"raspstore.github.io/file-manager/db"
	"raspstore.github.io/file-manager/repository"
)

func init() {

	err := godotenv.Load("../../.env")

	if err != nil {
		log.Panicln(err.Error())
	}
}

func TestUploadFile(t *testing.T) {
	cfg := db.NewConfig()

	serverAddress := fmt.Sprintf("localhost:%d", cfg.GrpcPort())
	flag.Parse()
	log.Printf("dial server %s", serverAddress)

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	client := pb.NewFileManagerServiceClient(conn)
	stream, err := client.Upload(context.Background())

	if err != nil {
		log.Panicln(err)
	}

	file, _ := os.Open("/home/murilobracero/test.txt")

	defer file.Close()

	req := &pb.CreateFileRequest{Data: &pb.CreateFileRequest_Filedata{
		Filedata: &pb.CreateFileRequestData{
			Filename:  "text.txt",
			CreatedBy: "test_user",
		},
	}}

	err = stream.Send(req)
	if err != nil {
		log.Panic(err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &pb.CreateFileRequest{
			Data: &pb.CreateFileRequest_Chunk{
				Chunk: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("file uploaded with id: %s, uri: %s", res.GetId(), res.GetUri())

}

func TestUpdateFile(t *testing.T) {
	cfg := db.NewConfig()
	c, err := db.NewMongoConnection(context.Background(), cfg)

	if err != nil {
		log.Panicln(err)
	}

	fr := repository.NewFilesRepository(context.Background(), c)

	files, err := fr.FindAll()

	if err != nil {
		log.Panicln(err)
	}

	uid := files[len(files)-1].Id.Hex()

	serverAddress := fmt.Sprintf("localhost:%d", cfg.GrpcPort())
	flag.Parse()
	log.Printf("dial server %s", serverAddress)

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	client := pb.NewFileManagerServiceClient(conn)
	stream, err := client.Update(context.Background())

	if err != nil {
		log.Panicln(err)
	}

	file, _ := os.Open("/home/murilobracero/test.txt")

	defer file.Close()

	req := &pb.UpdateFileRequest{Data: &pb.UpdateFileRequest_Filedata{
		Filedata: &pb.UpdateFileRequestData{
			Id:        uid,
			UpdatedBy: "test_user1",
		},
	}}

	err = stream.Send(req)
	if err != nil {
		log.Panic(err)
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &pb.UpdateFileRequest{
			Data: &pb.UpdateFileRequest_Chunk{
				Chunk: buffer[:n],
			},
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("file uploaded with id: %s, uri: %s", res.GetId(), res.GetUri())

}

func TestDeleteFile(t *testing.T) {
	cfg := db.NewConfig()
	c, err := db.NewMongoConnection(context.Background(), cfg)

	if err != nil {
		log.Panicln(err)
	}

	fr := repository.NewFilesRepository(context.Background(), c)

	files, err := fr.FindAll()

	if err != nil {
		log.Panicln(err)
	}

	serverAddress := fmt.Sprintf("localhost:%d", cfg.GrpcPort())
	flag.Parse()
	log.Printf("dial server %s", serverAddress)

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	uid := files[len(files)-1].Id.Hex()

	client := pb.NewFileManagerServiceClient(conn)
	res, err := client.Delete(context.Background(), &pb.GetFileRequest{Id: uid})

	if err != nil {
		log.Panicln(err)
	}

	assert.Equal(t, uid, res.Id)

}

func TestListFiles(t *testing.T) {
	cfg := db.NewConfig()

	serverAddress := fmt.Sprintf("localhost:%d", cfg.GrpcPort())
	flag.Parse()
	log.Printf("dial server %s", serverAddress)

	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	client := pb.NewFileManagerServiceClient(conn)
	stream, err := client.ListFiles(context.Background(), &pb.FindAllRequest{})

	if err != nil {
		log.Panicln(err)
	}

	for {

		res, err := stream.Recv()

		if err == io.EOF {

			log.Println("EOF")

			return

		}

		if err != nil {

			log.Printf("Err: %v", err)

		}

		log.Printf("output: %v", res)

	}

}
