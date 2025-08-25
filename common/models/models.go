package models

import (
	"time"

	"gorm.io/datatypes"
)

type Log struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Service   string         `json:"service"`
	Level     string         `json:"level"`
	Timestamp time.Time      `json:"timestamp"`
	Message   string         `json:"message"`
	Meta      datatypes.JSON `json:"meta"` // stored as JSONB
}

type AlertRule struct {
	ID        int      `gorm:"primaryKey" json:"id"`
	Name      string   `json:"name"`
	Service   string   `json:"service"` // optional
	Level     string   `json:"level"` // optional
	Keyword   string   `json:"keyword"` // match in message (optional)
	Field     string   `json:"field"` // JSONB field to inspect (e.g., latency)
	Condition string   `json:"condition"` // >, <, =, contains
	Threshold string   `json:"threshold"` // value to compare with
	Interval  string   `json:"interval"` // time window: "5m", "10m"
	Action    string   `json:"action"` // "log", "slack", "email", "webhook"
	Enabled   bool     `json:"enabled"`
}
