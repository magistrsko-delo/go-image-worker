package Worker

import (
	"encoding/json"
	"fmt"
	"go-image-worker/Http"
	"go-image-worker/Models"
	"go-image-worker/ffmpeg"
	"go-image-worker/grpc_client"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	pbTimeshift "go-image-worker/proto/timeshift_service"
	pbMediaMetadata "go-image-worker/proto/media_metadata"
	"strings"
	"time"
)

type Worker struct {
	RabbitMQ *RabbitMqConnection
	timeShiftGrpcClient *grpc_client.TimeShiftClient
	awsStorageGrpcClient *grpc_client.AwsStorageClient
	mediaMetadataGrpcClient *grpc_client.MediaMetadataClient
	mediaDownloader *Http.MediaDownloader
	ffmpeg *ffmpeg.FFmpeg
}

func (worker *Worker) Work()  {
	forever := make(chan bool)
	go func() {
		defer worker.timeShiftGrpcClient.Conn.Close()
		for d := range worker.RabbitMQ.msgs {
			log.Printf("Received a message: %s", d.Body)

			imageMediaRequest := &Models.ImageRabbitMqRequest{}
			err := json.Unmarshal(d.Body, imageMediaRequest)
			if err != nil{
				log.Println(err)
			}

			mediaChunksMetadata, err := worker.timeShiftGrpcClient.GetMediaChunkInfo(imageMediaRequest.MediaId)
			if err != nil{
				log.Println(err)
			}

			mediaMetadataInfo, err := worker.mediaMetadataGrpcClient.GetMediaMetadata(imageMediaRequest.MediaId)

			if err != nil{
				log.Println(err)
			}

			mediaChunks := mediaChunksMetadata.GetData()[0].GetChunks()  // 1080p chunks resolution
			tmpMediaLength := 0.0
			chunkIndex := -1
			chunkTime := -1.0 // only in seconds and mSeconds
			for i := 0; i < len(mediaChunks); i++ {
				if mediaChunks[i].GetLength() + tmpMediaLength >= imageMediaRequest.Time {
					chunkTime = imageMediaRequest.Time - tmpMediaLength
					chunkIndex = i
					break
				}

				tmpMediaLength = tmpMediaLength + mediaChunks[i].GetLength()
			}

			fmt.Println(chunkIndex)
			fmt.Println(chunkTime)

			err = worker.mediaDownloader.DownloadFile("./assets/chunks/" + mediaChunks[chunkIndex].GetAwsStorageName(), mediaChunks[chunkIndex].GetChunksUrl())
			if err != nil{
				log.Println(err)
			}

			err = worker.ffmpeg.ExecFFmpegCommand([]string{"-i", "assets/chunks/" + mediaChunks[chunkIndex].GetAwsStorageName(),
				"-acodec", "copy", "-vcodec", "copy", "assets/" + strings.Replace(mediaChunks[chunkIndex].GetAwsStorageName(), ".ts", ".mp4", 1)})

			imageThumbnailPath, err := worker.getMediaScreenShot(mediaChunks[chunkIndex], chunkTime, mediaMetadataInfo)
			fmt.Println("PATH THUMBNAIL: " + imageThumbnailPath)
			if err != nil{
				log.Println(err)
			}

			log.Println(imageThumbnailPath);
			mediaMetadataInfo.Thumbnail = imageThumbnailPath
			_, err = worker.mediaMetadataGrpcClient.UpdateMediaMetadata(mediaMetadataInfo)

			if err != nil{
				log.Println(err)
			}

			time.Sleep(1 * time.Second)

			// worker.removeFile("./assets/" + strings.Replace(mediaChunks[chunkIndex].GetAwsStorageName(), ".ts", ".mp4", 1))
			// worker.removeFile("./assets/chunks/" + mediaChunks[chunkIndex].GetAwsStorageName())

			log.Printf("Done")
			_ = d.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (worker *Worker) getMediaScreenShot(chunksData *pbTimeshift.ChunkResponse, chunkTime float64, mediaMetadata *pbMediaMetadata.MediaMetadataResponse) (string, error) {
	log.Println("CREATING MEDIA SCREENSHOT")

	var imageName string
	imageName = strconv.Itoa(int(chunksData.GetChunkId())) + "-" + strconv.Itoa(rand.Intn(1000000000000)) + "-" + strings.Replace(chunksData.GetAwsStorageName(), ".ts", ".jpg", 1)

	imageClipSec := strconv.Itoa(int(chunkTime))
	imageClipMiliSec := strconv.Itoa(int( math.Mod(chunkTime * 1000, 1000) ))
	ss := "00:00:" + imageClipSec + "." + imageClipMiliSec
	err := worker.ffmpeg.ExecFFmpegCommand([]string{"-ss", ss , "-i",
		"assets/" + strings.Replace(chunksData.GetAwsStorageName(), ".ts", ".mp4", 1),
		"-vframes", "1", "-g:v", "2", "./assets/" + imageName})
	if err != nil {
		log.Println(err)
		return "" , err
	}
	_, err = worker.awsStorageGrpcClient.UploadMedia( "./assets/" + imageName, "mag20-images", imageName)  // TODO for later add this to configuration maybe..
	if err != nil {
		log.Println(err)
		return "" , err
	}

	time.Sleep(2 * time.Second)

	// worker.removeFile("./assets/" + imageName)
	return "v1/mediaManager/mag20-images/" + imageName, nil
}

func (worker *Worker) removeFile(path string)  {
	err := os.Remove(path)
	if err != nil {
		fmt.Println(err)
	}
}

func InitWorker() *Worker {
	return &Worker{
		RabbitMQ:            		initRabbitMqConnection(Models.GetEnvStruct()),
		timeShiftGrpcClient: 		grpc_client.InitTimeShiftClient(),
		awsStorageGrpcClient: 		grpc_client.InitAwsStorageGrpcClient(),
		mediaMetadataGrpcClient: 	grpc_client.InitMediaMetadataGrpcClient(),
		mediaDownloader: 			&Http.MediaDownloader{},
		ffmpeg: 					&ffmpeg.FFmpeg{},
	}
}
