package grpc_client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"main/Models"
	pbMediaMetadata "main/proto/media_metadata"
	"strconv"
)


type MediaMetadataClient struct {
	Conn *grpc.ClientConn
	client pbMediaMetadata.MediaMetadataClient
}

func (mediaMetadataClient *MediaMetadataClient) UpdateMediaMetadata(mediaMetadata *Models.MediaMetadata) (*pbMediaMetadata.MediaMetadataResponse, error)  {

	createdAt, _ := strconv.ParseInt(mediaMetadata.CreatedAt, 10, 64)

	response, err := mediaMetadataClient.client.UpdateMediaMetadata(context.Background(), &pbMediaMetadata.UpdateMediaRequest{
		MediaId:                  int32(mediaMetadata.MediaId),
		Name:                     mediaMetadata.Name,
		SiteName:                 mediaMetadata.SiteName,
		Length:                   int32(mediaMetadata.Length),
		Status:                   3,
		Thumbnail:                mediaMetadata.Thumbnail,
		ProjectId:                int32(mediaMetadata.ProjectId),
		AwsBucketWholeMedia:      mediaMetadata.AwsBucketWholeMedia,
		AwsStorageNameWholeMedia: mediaMetadata.AwsStorageNameWholeMedia,
		CreatedAt:                createdAt,
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}


func InitMediaMetadataGrpcClient() *MediaMetadataClient  {
	env := Models.GetEnvStruct()
	fmt.Println("CONNECTING")
	conn, err := grpc.Dial(env.MediaMetadataGrpcServer + ":" + env.MediaMetadataGrpcPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	fmt.Println("END CONNECTION")

	client := pbMediaMetadata.NewMediaMetadataClient(conn)
	return &MediaMetadataClient{
		Conn:    conn,
		client:  client,
	}

}
