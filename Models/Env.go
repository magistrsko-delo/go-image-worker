package Models

import (
	"fmt"
	"os"
)

var envStruct *Env

type Env struct {
	RabbitUser string
	RabbitPassword string
	RabbitQueue string
	RabbitHost string
	RabbitPort string
	Env string
	AwsStorageGrpcServer string
	AwsStorageGrpcPort string
	TimeShiftGrpcServer string
	TimeShiftGrpcPort string
}

func InitEnv()  {
	envStruct = &Env{
		RabbitUser:       			os.Getenv("RABBIT_USER"),
		RabbitPassword:   			os.Getenv("RABBIT_PASSWORD"),
		RabbitQueue:      			os.Getenv("RABBIT_QUEUE"),
		RabbitHost:       			os.Getenv("RABBIT_HOST"),
		RabbitPort: 				os.Getenv("RABBIT_PORT"),
		Env: 			  			os.Getenv("ENV"),
		AwsStorageGrpcServer: 		os.Getenv("AWS_STORAGE_GRPC_SERVER"),
		AwsStorageGrpcPort:			os.Getenv("AWS_STORAGE_GRPC_PORT"),
		TimeShiftGrpcServer:        os.Getenv("TIMESHIFT_GRPC_SERVER"),
		TimeShiftGrpcPort: 			os.Getenv("TIMESHIFT_GRPC_PORT"),
	}
	fmt.Println(envStruct)
}

func GetEnvStruct() *Env  {
	return  envStruct
}