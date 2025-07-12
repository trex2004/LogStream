package main

import (
	"log"

	// "github.com/nats-io/nats.go"	
	"github.com/trex2004/logstream/common/db"
)

func main(){

	// nc, err := nats.Connect(nats.DefaultURL)
    // if err != nil {
    //     log.Fatalf("Error connecting to NATS: %v", err)
    // }
	// js, err := nc.JetStream()
	// if err != nil {
	// 	log.Fatalf("Error creating JetStream context: %v", err)
	// }

	LogStoreDB, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer LogStoreDB.Close()
}