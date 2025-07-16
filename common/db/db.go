package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/trex2004/logstream/common/models"
)

type LogStoreDB struct {
	DB *sql.DB
}

func InitDB() (*LogStoreDB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	
	url := os.Getenv("POSTGRES_URL")
	log.Printf("Connecting to database with URL: %s", url)
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Connected to the database successfully")
	return &LogStoreDB{DB: db}, nil
}

func (l *LogStoreDB) CreateLogStoreTable() error {
	if _,err := l.DB.Exec(`
		CREATE TABLE IF NOT EXISTS logs (
			id SERIAL PRIMARY KEY,
			service TEXT NOT NULL,
			level TEXT NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL,
			message TEXT NOT NULL,
			meta JSONB
		);
	`); err != nil {
		return err
	}
	log.Println("Log store table created successfully")
	return nil
}

func (l *LogStoreDB) Close() error {
	if err := l.DB.Close(); err != nil {
		return err
	}
	log.Println("Database connection closed successfully")
	return nil
}

func InsertLogMessage(db *LogStoreDB, logMsg models.Log) error {
	_, err := db.DB.Exec(`
		INSERT INTO logs (service, level, timestamp, message, meta)
		VALUES ($1, $2, $3, $4, $5)
	`, logMsg.Service, logMsg.Level, logMsg.Timestamp, logMsg.Message, toJSONB(logMsg.Meta))
	if err != nil {
		return err
	}
	log.Println("Log message inserted successfully")
	return nil
}

func toJSONB(m map[string]interface{}) []byte {
	data, _ := json.Marshal(m)
	return data
}

func (l *LogStoreDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := l.DB.Query(query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	log.Printf("Query executed successfully: %s", query)
	return rows, nil
}