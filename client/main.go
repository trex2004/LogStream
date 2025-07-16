package main

import (
	"context"
	"log"
	"time"

	"github.com/trex2004/logstream/common"
	pb "github.com/trex2004/logstream/proto"
	"google.golang.org/grpc"
)

var (
	grpcPort = common.GetEnv("GRPC_PORT", ":50051")
)

func main(){
	conn,err := grpc.Dial(grpcPort,grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()
	log.Printf("Connected to gRPC server on %s", grpcPort)

	client := pb.NewLogServiceClient(conn)
	log.Printf("Sending test log to gRPC server...")

	for {
		time.Sleep(5*time.Second)
		res,err := client.SendLog(context.Background(), &pb.LogRequest{
			Service:   "auth-service",
			Level:     "ERROR",
			Timestamp: time.Now().Format(time.RFC3339),
			Message:   "This is a test log message",
			Meta:  map[string]string{"user_name": "trex2004"},
		})
		if err != nil {
			log.Fatalf("Error sending log: %v", err)
		}
		log.Printf("Log response: %+v", res)
	}
}