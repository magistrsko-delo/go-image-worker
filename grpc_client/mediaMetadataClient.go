package grpc_client

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"go-image-worker/Models"
	pbMediaMetadata "go-image-worker/proto/media_metadata"
)


type MediaMetadataClient struct {
	Conn *grpc.ClientConn
	client pbMediaMetadata.MediaMetadataClient
}

func (mediaMetadataClient *MediaMetadataClient) UpdateMediaMetadata(mediaMetadata *pbMediaMetadata.MediaMetadataResponse) (*pbMediaMetadata.MediaMetadataResponse, error)  {

	response, err := mediaMetadataClient.client.UpdateMediaMetadata(context.Background(), &pbMediaMetadata.UpdateMediaRequest{
		MediaId:                  mediaMetadata.GetMediaId(),
		Name:                     mediaMetadata.GetName(),
		SiteName:                 mediaMetadata.GetSiteName(),
		Length:                   mediaMetadata.GetLength(),
		Status:                   mediaMetadata.GetStatus(),
		Thumbnail:                mediaMetadata.GetThumbnail(),
		ProjectId:                mediaMetadata.GetProjectId(),
		AwsBucketWholeMedia:      mediaMetadata.GetAwsBucketWholeMedia(),
		AwsStorageNameWholeMedia: mediaMetadata.GetAwsStorageNameWholeMedia(),
		CreatedAt:                mediaMetadata.GetCreatedAt(),
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (mediaMetadataClient *MediaMetadataClient) GetMediaMetadata(mediaId int32) (*pbMediaMetadata.MediaMetadataResponse, error) {

	response, err := mediaMetadataClient.client.GetMediaMetadata(context.Background(), &pbMediaMetadata.GetMediaMetadataRequest{
		MediaId:              mediaId,
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
