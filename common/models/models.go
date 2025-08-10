package models

import "time"

type Log struct {
	Service  string          `json:"service"`
	Level    string          `json:"level"`
	Timestamp time.Time          `json:"timestamp"`
	Message  string          `json:"message"`
	Meta     map[string]interface{} `json:"meta"`
}

type AlertRule struct {
	ID        int
	Name      string  `json:"name"`
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
