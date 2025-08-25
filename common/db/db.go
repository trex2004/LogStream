package db

import (
	"encoding/json"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/trex2004/logstream/common/models"
)

type LogStoreDB struct {
	DB *gorm.DB
}

func InitDB() (*LogStoreDB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	
	url := os.Getenv("POSTGRES_URL")
	log.Printf("Connecting to database with URL: %s", url)
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&models.Log{}, &models.AlertRule{}); err != nil {
		return nil, err
	}
	log.Println("Connected to the database successfully")
	return &LogStoreDB{DB: db}, nil
}

func (l *LogStoreDB) CreateLogStoreTable() error {
	return l.DB.AutoMigrate(&models.Log{})
}

func (l *LogStoreDB) CreateAlertRulesTable() error {
	return l.DB.AutoMigrate(&models.AlertRule{})
}

func (l *LogStoreDB) Close() error {
	sqlDB, err := l.DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

//what is this?? change this next ,, what have i even done??
func InsertLogMessage(db *LogStoreDB, logMsg models.Log) error {
	if err := db.DB.Create(&logMsg).Error; err != nil {
		return err
	}
	log.Println("Log message inserted successfully")
	return nil
}

func toJSONB(m map[string]interface{}) []byte {
	data, _ := json.Marshal(m)
	return data
}

// func (l *LogStoreDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
// 	rows, err := l.DB.Query(query, args...)
// 	if err != nil {
// 		log.Printf("Error executing query: %v", err)
// 		return nil, err
// 	}
// 	log.Printf("Query executed successfully: %s", query)
// 	return rows, nil
// }

// func (l *LogStoreDB) QueryRow(query string, args ...interface{}) *sql.Row {
// 	return l.DB.QueryRow(query, args...)
// }

func (l *LogStoreDB) InsertAlertRule(alert models.AlertRule) error {
	if err := l.DB.Create(&alert).Error; err != nil {
		return err
	}
	log.Println("Alert rule inserted successfully")
	return nil
}