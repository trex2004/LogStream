package main

import (
	// "encoding/json"
	"encoding/json"
	"log"
	"strconv"

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

		service := c.Query("service")
		level := c.Query("level")
		limit := c.DefaultQuery("limit", "100")

		query := "SELECT service, level, timestamp, message, meta FROM logs WHERE 1=1"
		args := []interface{}{}
		i := 1

		if service != "" {
			query += " AND service = $" + strconv.Itoa(i)
			args = append(args, service)
			i++
		}

		if level != "" {
			query += " AND level = $" + strconv.Itoa(i)
			args = append(args, level)
			i++
		}

		query += " ORDER BY timestamp DESC"

		query += " LIMIT $" + strconv.Itoa(i)
		args = append(args, limit)
		i++

		log.Printf("Executing query: %s with args: %v", query, args)

		rows, err := LogStoreDB.Query(query, args...)
		if err != nil {
			log.Printf("Error querying logs: %v", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()

		logs := make([]map[string]interface{}, 0)
		for rows.Next() {
			log.Printf("Processing row...")
			var service, level, message string
			var timestamp string
			var meta map[string]interface{}
			var metaRaw []byte

			err := rows.Scan(&service, &level, &timestamp, &message, &metaRaw)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}

			json.Unmarshal(metaRaw, &meta)

			logs = append(logs, gin.H{
				"service":   service,
				"level":     level,
				"timestamp": timestamp,
				"message":   message,
				"meta":      meta,
			})
    	}

		c.JSON(200, logs)
		
	})

	log.Printf("LogStream API is running on port 8080")
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}