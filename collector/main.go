package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/trex2004/logstream/common"
	pb "github.com/trex2004/logstream/proto"

	"google.golang.org/grpc"
)

var (
	grpcPort = common.GetEnv("GRPC_PORT", ":50051")
)

type server struct {
    pb.UnimplementedLogServiceServer
    nc *nats.Conn
}

func (s *server) SendLog(ctx context.Context, req *pb.LogRequest) (*pb.LogResponse, error) {
    log.Printf("Received log from %s: %s", req.Service, req.Message)

    meta,err := json.Marshal(req.Meta)
    if err != nil {
        log.Printf("Error marshalling meta: %v", err)
        return &pb.LogResponse{Success: false, Message: "Failed to marshal meta"}, err
    }


    logMsg := fmt.Sprintf(`{"service":"%s","level":"%s","timestamp":"%s","message":"%s","meta":%s}`, req.Service, req.Level, req.Timestamp, req.Message, meta)

    log.Printf("Publishing log message to NATS: %s", logMsg)

    subject := fmt.Sprintf("logs.%s.%s", req.Service, req.Level)
    err = s.nc.Publish(subject, []byte(logMsg))
	// log.Printf("Error publishing to NATS: %v", err)
    if err != nil {
        return &pb.LogResponse{Success: false, Message: "Failed to publish to NATS"}, err
    }

    return &pb.LogResponse{Success: true, Message: "Log received"}, nil
}

func main() {
    nc, err := nats.Connect(os.Getenv("NATS_URL"))
    if err != nil {
        log.Fatalf("Error connecting to NATS: %v", err)
    }
	err = common.InitiateStreams(nc)
	if err != nil {
		log.Fatalf("Error initiating streams: %v", err)
	}

    lis, err := net.Listen("tcp", grpcPort)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }
	defer lis.Close()

    s := grpc.NewServer()
    pb.RegisterLogServiceServer(s, &server{nc: nc})

    log.Printf("gRPC Server listening on %s", grpcPort)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
