package grpc_client

import (
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"main/Models"
	pbAwsStorage "main/proto/aws_storage"
	"os"
)

type AwsStorageClient struct {
	Conn *grpc.ClientConn
	client pbAwsStorage.AwsstorageClient
}

func (awsStorageClient *AwsStorageClient) UploadMedia(filePath string, bucketName string, storageName string) (*pbAwsStorage.UploadResponse, error)  {
	fmt.Println("UPLOADING: " + storageName)
	var (
		writing = true
		buf     []byte
		n       int
		file    *os.File
	)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stream, err := awsStorageClient.client.UploadFile(context.Background())

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	defer stream.CloseSend()

	buf = make([]byte,  512 * 1024) // 512k

	for writing {
		n, err = file.Read(buf)
		if err != nil {
			if err == io.EOF {
				writing = false
				err = nil
				continue
			}
			err = errors.New("errored while copying from file to buffer")
			return nil, err
		}

		err = stream.Send(&pbAwsStorage.UploadRequest{
			Bucketname:           bucketName,
			Medianame:            storageName,
			Data:                 buf[:n],
		})

		if err != nil {
			err = errors.New("failed to send chunk via stream")
			return nil, err
		}
	}

	resp, err := stream.CloseAndRecv()

	if err != nil {
		err = errors.New("failed to receive upstream status response")
		return nil, err
	}

	return resp, nil
}

func InitAwsStorageGrpcClient() *AwsStorageClient {

	env := Models.GetEnvStruct()
	fmt.Println("CONNECTING aws storageClient")
	conn, err := grpc.Dial(env.AwsStorageGrpcServer + ":" + env.AwsStorageGrpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	fmt.Println("END CONNECTION aws storage client")

	client := pbAwsStorage.NewAwsstorageClient(conn)
	return &AwsStorageClient{
		Conn:    conn,
		client:  client,
	}

}



