package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/trex2004/logstream/common/db"
	"github.com/trex2004/logstream/common/models"
)



func main(){

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
    if err != nil {
        log.Fatalf("Error connecting to NATS: %v", err)
    }
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error creating JetStream context: %v", err)
	}

	LogStoreDB, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer LogStoreDB.Close()

	err = LogStoreDB.CreateLogStoreTable()
	if err != nil {
		log.Fatalf("Error creating log store table: %v", err)
	}

	err = LogStoreDB.CreateAlertRulesTable()
	if err != nil {
		log.Fatalf("Error creating alert rule table: %v", err)
	}

	_,err = js.Subscribe("logs.*.*", func(msg *nats.Msg) {
		var logMsg models.Log
		err := json.Unmarshal(msg.Data, &logMsg)
		if err != nil {
			log.Printf("Error unmarshalling log message: %v", err)
			msg.Ack()
			return
		}

		err = db.InsertLogMessage(LogStoreDB, logMsg)
		if err != nil {
			log.Printf("Error inserting log message into database: %v", err)
			msg.Nak()
			return
		} 
		log.Printf("Log message inserted into database: %s", logMsg.Message)
		msg.Ack()
	},nats.Durable("logstream_durable_consumer"), nats.ManualAck())

	if err != nil {
		log.Fatalf("Error subscribing to NATS subject: %v", err)
	}
	log.Println("LogStream Processor is running and listening for log messages...")

	select {}
}