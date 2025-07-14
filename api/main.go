package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/trex2004/logstream/common/db"
)

func main(){

	LogStoreDB, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer LogStoreDB.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(200, "LogStream API is running")
	})

	router.GET("/logs", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "LogStream API is running",
		})
	})

	log.Printf("LogStream API is running on port 8080")
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}