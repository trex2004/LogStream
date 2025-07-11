package common

import (
	"log"

	"github.com/nats-io/nats.go"
)

func InitiateStreams(nc *nats.Conn) (error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil
	}

	//check if the stream already exists
	stream, err := js.StreamInfo("logs")
	if err == nil {
		log.Printf("Stream 'logs' already exists: %+v", stream)
		return nil
	}

	log.Printf("Creating JetStream streams...")

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "logs",
		Subjects: []string{"logs.*.*"},
		Storage: nats.FileStorage,
		Retention: nats.WorkQueuePolicy,
	})
	if err != nil{
		return err
	}
	log.Println("JetStream streams created successfully")
	return nil
}