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

		rows, err := LogStoreDB.DB.Query(query, args...)
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

	router.GET("/metrics/count", func(c *gin.Context){
		service := c.Query("service")
		level := c.Query("level")

		query := "SELECT COUNT(*) FROM logs WHERE 1=1"
		args := []interface{}{}
		i := 1

		if service != ""{
			query += " AND service = $" + strconv.Itoa(i)
			args = append(args, service)
			i++
		}
		
		if level != ""{
			query += " AND level = $" + strconv.Itoa(i)
			args = append(args, level)
			i++
		}

		log.Printf("Executing count query: %s with args: %v", query, args)
		var count int
		err := LogStoreDB.QueryRow(query, args...).Scan(&count)
		if err != nil {
			log.Printf("Error querying count: %v", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		c.JSON(200, gin.H{"count": count})
	})

	router.GET("metrics/errors-per-service", func(c *gin.Context) {
		query := `SELECT service, COUNT(*) as error_count FROM logs WHERE level = 'ERROR' GROUP BY service`
		rows, err := LogStoreDB.Query(query)
		if err != nil {
			log.Printf("Error querying errors per service: %v", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()

		errors := make(map[string]int)
		for rows.Next() {
			var service string
			var errorCount int
			if err := rows.Scan(&service, &errorCount); err != nil {
				log.Printf("Error scanning row: %v", err)
				c.JSON(500, gin.H{"error": "Internal Server Error"})
				return
			}
			errors[service] = errorCount
		}
		c.JSON(200, errors)
	})

	router.GET("/metrics/daily-logs", func(c *gin.Context){
		daysStr := c.DefaultQuery("days", "7")
		days, err := strconv.Atoi(daysStr)
		if err != nil || days <= 0 {
			log.Printf("Invalid 'days' value: %s", daysStr)
			c.JSON(400, gin.H{"error": "Invalid 'days' parameter"})
			return
		}
		query := `SELECT DATE(timestamp) , COUNT(*) FROM logs WHERE timestamp > NOW() - INTERVAL '` + strconv.Itoa(days)  + ` days' GROUP BY DATE(timestamp) ORDER BY DATE(timestamp) DESC`
		log.Printf("Executing daily logs query: %s ", query)

		rows, err := LogStoreDB.Query(query)
		if err != nil {
			log.Printf("Error querying daily logs: %v", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()

		dailyLogs := make(map[string]int)
		for rows.Next() {
			var date string
			var logCount int
			if err := rows.Scan(&date, &logCount); err != nil {
				log.Printf("Error scanning row: %v", err)
				c.JSON(500, gin.H{"error": "Internal Server Error"})
				return
			}
			dailyLogs[date]=logCount
		}
		
		c.JSON(200, dailyLogs)
	})
	
	router.GET("metrics/levels-distribution", func(c *gin.Context){
		daysStr := c.DefaultQuery("days", "7")
		days,err := strconv.Atoi(daysStr)	
		if err!=nil {
			log.Printf("Invalid days %s",daysStr)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		
		query := `SELECT level, COUNT(*) FROM logs WHERE timestamp>NOW() - INTERVAL '`+ strconv.Itoa(days) +` days' GROUP BY level`
		log.Printf("Executing daily level logs query: %s ", query)
		
		rows,err := LogStoreDB.Query(query)
		if err!=nil{
			log.Printf("Error querying daily level logs: %v", err)
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}
		defer rows.Close()
		
		levelLogs := make(map[string]int)
		for rows.Next() {
			var level string
			var count int
			if err := rows.Scan(&level, &count); err!=nil {
				log.Printf("Error scanning row: %v", err)
				c.JSON(500, gin.H{"error": "Internal Server Error"})
				return
			}
			levelLogs[level]=count
		}

		c.JSON(200, levelLogs)

	})
	
	log.Printf("LogStream API is running on port 8080")
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}