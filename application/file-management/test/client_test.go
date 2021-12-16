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
	"google.golang.org/grpc"
	"raspstore.github.io/file-manager/db"
	"raspstore.github.io/file-manager/pb"
)

func init() {
	err := godotenv.Load("../.env")

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
			CreatedBy: "zephiroca",
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
