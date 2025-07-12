package db

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type LogStoreDB struct {
	DB *sql.DB
}

func InitDB() (*LogStoreDB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	
	url := os.Getenv("POSTGRESQL_URL")
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