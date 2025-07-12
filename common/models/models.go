package models

import "time"

type Log struct {
	Service  string          `json:"service"`
	Level    string          `json:"level"`
	Timestamp time.Time          `json:"timestamp"`
	Message  string          `json:"message"`
	Meta     map[string]interface{} `json:"meta"`
}