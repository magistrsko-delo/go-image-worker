package Worker

import (
	"go-image-worker/Models"
	"go-image-worker/grpc_client"
	"log"
)

type Worker struct {
	RabbitMQ *RabbitMqConnection
	timeShiftGrpcClient *grpc_client.TimeShiftClient
}

func (worker *Worker) Work()  {
	forever := make(chan bool)
	go func() {
		defer worker.timeShiftGrpcClient.Conn.Close()
		for d := range worker.RabbitMQ.msgs {
			log.Printf("Received a message: %s", d.Body)

			log.Printf("Done")
			_ = d.Ack(false)
		}
	}()
	<-forever
}

func InitWorker() *Worker {
	return &Worker{
		RabbitMQ:            initRabbitMqConnection(Models.GetEnvStruct()),
		timeShiftGrpcClient: grpc_client.InitTimeShiftClient(),
	}
}
