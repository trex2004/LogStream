package main

import (
	// "encoding/json"
	// "encoding/csv"
	// "encoding/json"
	// "fmt"
	"log"
	"strconv"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/trex2004/logstream/common/db"
	"github.com/trex2004/logstream/common/models"
)

func main() {

	LogStoreDB, err := db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer LogStoreDB.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(200, "LogStream API is running")
	})

	// make log into a model in common
	router.GET("/logs", func(c *gin.Context) {
		service := c.Query("service")
		level := c.Query("level")
		from := c.Query("from")
		to := c.Query("to")

		limitStr := c.DefaultQuery("limit", "100")
		offsetStr := c.DefaultQuery("offset", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			c.JSON(400, gin.H{"error": "Invalid 'limit' parameter"})
			return
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			c.JSON(400, gin.H{"error": "Invalid 'offset' parameter"})
			return
		}

		var logs []models.Log
		db := LogStoreDB.DB

		if service != "" {
			db = db.Where("service = ?", service)
		}
		if level != "" {
			db = db.Where("level = ?", level)
		}
		if from != "" {
			db = db.Where("timestamp >= ?", from)
		}
		if to != "" {
			db = db.Where("timestamp <= ?", to)
		}

		err = db.Order("timestamp DESC").
			Limit(limit).
			Offset(offset).
			Find(&logs).Error
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal Server Error"})
			return
		}

		// Meta is already a map[string]interface{}, no extra work
		c.JSON(200, logs)
	})

	// router.GET("/metrics/count", func(c *gin.Context) {
	// 	service := c.Query("service")
	// 	level := c.Query("level")

	// 	query := "SELECT COUNT(*) FROM logs WHERE 1=1"
	// 	args := []interface{}{}
	// 	i := 1

	// 	if service != "" {
	// 		query += " AND service = $" + strconv.Itoa(i)
	// 		args = append(args, service)
	// 		i++
	// 	}

	// 	if level != "" {
	// 		query += " AND level = $" + strconv.Itoa(i)
	// 		args = append(args, level)
	// 		i++
	// 	}

	// 	log.Printf("Executing count query: %s with args: %v", query, args)
	// 	var count int
	// 	err := LogStoreDB.QueryRow(query, args...).Scan(&count)
	// 	if err != nil {
	// 		log.Printf("Error querying count: %v", err)
	// 		c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 		return
	// 	}
	// 	c.JSON(200, gin.H{"count": count})
	// })

	// router.GET("metrics/errors-per-service", func(c *gin.Context) {
	// 	query := `SELECT service, COUNT(*) as error_count FROM logs WHERE level = 'ERROR' GROUP BY service`
	// 	rows, err := LogStoreDB.Query(query)
	// 	if err != nil {
	// 		log.Printf("Error querying errors per service: %v", err)
	// 		c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 		return
	// 	}
	// 	defer rows.Close()

	// 	errors := make(map[string]int)
	// 	for rows.Next() {
	// 		var service string
	// 		var errorCount int
	// 		if err := rows.Scan(&service, &errorCount); err != nil {
	// 			log.Printf("Error scanning row: %v", err)
	// 			c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 			return
	// 		}
	// 		errors[service] = errorCount
	// 	}
	// 	c.JSON(200, errors)
	// })

	// router.GET("/metrics/daily-logs", func(c *gin.Context) {
	// 	daysStr := c.DefaultQuery("days", "7")
	// 	days, err := strconv.Atoi(daysStr)
	// 	if err != nil || days <= 0 {
	// 		log.Printf("Invalid 'days' value: %s", daysStr)
	// 		c.JSON(400, gin.H{"error": "Invalid 'days' parameter"})
	// 		return
	// 	}
	// 	query := `SELECT DATE(timestamp) , COUNT(*) FROM logs WHERE timestamp > NOW() - INTERVAL '` + strconv.Itoa(days) + ` days' GROUP BY DATE(timestamp) ORDER BY DATE(timestamp) DESC`
	// 	log.Printf("Executing daily logs query: %s ", query)

	// 	rows, err := LogStoreDB.Query(query)
	// 	if err != nil {
	// 		log.Printf("Error querying daily logs: %v", err)
	// 		c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 		return
	// 	}
	// 	defer rows.Close()

	// 	dailyLogs := make(map[string]int)
	// 	for rows.Next() {
	// 		var date string
	// 		var logCount int
	// 		if err := rows.Scan(&date, &logCount); err != nil {
	// 			log.Printf("Error scanning row: %v", err)
	// 			c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 			return
	// 		}
	// 		dailyLogs[date] = logCount
	// 	}

	// 	c.JSON(200, dailyLogs)
	// })

	// router.GET("metrics/levels-distribution", func(c *gin.Context) {
	// 	daysStr := c.DefaultQuery("days", "7")
	// 	days, err := strconv.Atoi(daysStr)
	// 	if err != nil || days <= 0 {
	// 		log.Printf("Invalid days %s", daysStr)
	// 		c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 		return
	// 	}

	// 	query := `SELECT level, COUNT(*) FROM logs WHERE timestamp>NOW() - INTERVAL '` + strconv.Itoa(days) + ` days' GROUP BY level`
	// 	log.Printf("Executing level logs query: %s ", query)

	// 	rows, err := LogStoreDB.Query(query)
	// 	if err != nil {
	// 		log.Printf("Error querying level logs: %v", err)
	// 		c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 		return
	// 	}
	// 	defer rows.Close()

	// 	levelLogs := make(map[string]int)
	// 	for rows.Next() {
	// 		var level string
	// 		var count int
	// 		if err := rows.Scan(&level, &count); err != nil {
	// 			log.Printf("Error scanning row: %v", err)
	// 			c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 			return
	// 		}
	// 		levelLogs[level] = count
	// 	}

	// 	c.JSON(200, levelLogs)

	// })

	// router.GET("/metrics/services-activity", func(c *gin.Context) {
	// 	daysStr := c.DefaultQuery("days", "7")
	// 	days, err := strconv.Atoi(daysStr)
	// 	if err != nil || days <= 0 {
	// 		log.Printf("Invalid Days")
	// 		c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 		return
	// 	}

	// 	query := `SELECT service, COUNT(*) FROM logs WHERE timestamp>NOW() - INTERVAL '` + strconv.Itoa(days) + ` days' GROUP BY service ORDER BY COUNT(*) DESC`
	// 	log.Printf("Executing service logs query: %s ", query)

	// 	rows, err := LogStoreDB.Query(query)
	// 	if err != nil {
	// 		log.Printf("Error service level logs: %v", err)
	// 		c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 		return
	// 	}
	// 	defer rows.Close()

	// 	serviceLogs := make(map[string]int)
	// 	for rows.Next() {
	// 		var service string
	// 		var count int
	// 		if err := rows.Scan(&service, &count); err != nil {
	// 			log.Printf("Error scanning row: %v", err)
	// 			c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 			return
	// 		}
	// 		serviceLogs[service] = count
	// 	}

	// 	c.JSON(200, serviceLogs)
	// })

	// router.GET("/logs/export", func(c *gin.Context) {
	// 	service := c.Query("service")
	// 	level := c.Query("level")
	// 	from := c.Query("from")
	// 	to := c.Query("to")
	// 	limitStr := c.DefaultQuery("limit", "100")

	// 	limit, err := strconv.Atoi(limitStr)
	// 	if err != nil || limit <= 0 {
	// 		log.Printf("Invalid 'limit' value: %s", limitStr)
	// 		c.JSON(400, gin.H{"error": "Invalid 'limit' parameter"})
	// 		return
	// 	}

	// 	query := "SELECT service, level, timestamp, message, meta FROM logs WHERE 1=1"
	// 	args := []interface{}{}
	// 	i := 1

	// 	if service != "" {
	// 		query += fmt.Sprintf(" AND service = $%d", i)
	// 		args = append(args, service)
	// 		i++
	// 	}
	// 	if level != "" {
	// 		query += fmt.Sprintf(" AND level = $%d", i)
	// 		args = append(args, level)
	// 		i++
	// 	}
	// 	if from != "" {
	// 		query += fmt.Sprintf(" AND timestamp >= $%d", i)
	// 		args = append(args, from)
	// 		i++
	// 	}
	// 	if to != "" {
	// 		query += fmt.Sprintf(" AND timestamp <= $%d", i)
	// 		args = append(args, to)
	// 		i++
	// 	}

	// 	query += fmt.Sprintf(" ORDER BY timestamp DESC LIMIT $%d", i)
	// 	args = append(args, limit)

	// 	log.Printf("Executing export query: %s with args: %v", query, args)

	// 	rows, err := LogStoreDB.Query(query, args...)
	// 	if err != nil {
	// 		log.Printf("Error querying logs for export: %v", err)
	// 		c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 		return
	// 	}
	// 	defer rows.Close()

	// 	c.Header("Content-Disposition", "attachment; filename=logs.csv")
	// 	c.Header("Content-Type", "text/csv")

	// 	writer := csv.NewWriter(c.Writer)
	// 	defer writer.Flush()

	// 	writer.Write([]string{"Service", "Level", "Timestamp", "Message", "Meta"})

	// 	for rows.Next() {
	// 		var service, level, message string
	// 		var timestamp time.Time
	// 		var metaRaw []byte

	// 		if err := rows.Scan(&service, &level, &timestamp, &message, &metaRaw); err != nil {
	// 			log.Printf("Error scanning row for export: %v", err)
	// 			c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 			return
	// 		}

	// 		meta := make(map[string]interface{})
	// 		if err := json.Unmarshal(metaRaw, &meta); err != nil {
	// 			log.Printf("Error decoding meta for export: %v", err)
	// 			meta = nil
	// 		}

	// 		log.Printf("Meta for export: %+v", meta)

	// 		writer.Write([]string{
	// 			service,
	// 			level,
	// 			timestamp.Format(time.RFC3339),
	// 			message,
	// 			string(metaRaw), // Write raw JSON for simplicity
	// 		})
	// 	}

	// })

	// router.POST("/alerts", func(c *gin.Context) {
	// 	var alert models.AlertRule

	// 	if err := c.BindJSON(&alert); err != nil {
	// 		log.Printf("Error binding alert rule: %v", err)
	// 		c.JSON(400, gin.H{"error": "Invalid request body"})
	// 		return
	// 	}
	// 	//add more verificcation here
	// 	if alert.Name == "" || alert.Action == "" {
	// 		log.Printf("Missing required fields in alert rule")
	// 		c.JSON(400, gin.H{"error": "Name and Action are required fields"})
	// 		return
	// 	}
	// 	if alert.Interval == "" {
	// 		alert.Interval = "5m"
	// 	}
	// 	if err := LogStoreDB.InsertAlertRule(alert); err != nil {
	// 		log.Printf("Error inserting alert rule: %v", err)
	// 		c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 		return
	// 	}
	// 	log.Printf("Alert rule created successfully: %+v", alert)
	// 	c.JSON(201, gin.H{"message": "Alert rule created successfully", "alert": alert})
	// })

	// router.GET("/alerts", func(c *gin.Context) {
	// 	query := "Select * from alert_rules"

	// 	rows,err := LogStoreDB.Query(query)
	// 	if err!=nil {
	// 		log.Printf("Error querying for logs...")
	// 		c.JSON(500,gin.H{"error":"Internal Server error"})
	// 		return
	// 	}
	// 	defer rows.Close()

	// 	alerts := make([]models.AlertRule,0)
	// 	for rows.Next() {
	// 		var alert models.AlertRule
	// 		if err := rows.Scan(&alert.ID, &alert.Name, &alert.Service, &alert.Level, &alert.Keyword, &alert.Field, &alert.Condition, &alert.Threshold, &alert.Interval, &alert.Action, &alert.Enabled); err != nil {
	// 			log.Printf("Error scanning alert rule: %v", err)
	// 			c.JSON(500, gin.H{"error": "Internal Server Error"})
	// 			return
	// 		}
	// 		alerts = append(alerts, alert)
	// 	}
	// 	c.JSON(200, alerts)
	// })

	log.Printf("LogStream API is running on port 8080")
	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
