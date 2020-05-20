package main

import (
	"crypto/tls"
	"github.com/joho/godotenv"
	"go-image-worker/Models"
	"go-image-worker/Worker"
	"log"
	"net/http"
)

func init()  {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	Models.InitEnv()
}


func main()  {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	worker := Worker.InitWorker()
	defer worker.RabbitMQ.Conn.Close()
	defer worker.RabbitMQ.Ch.Close()
	worker.Work()
}
