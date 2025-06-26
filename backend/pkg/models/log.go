package models

import (
	"database/sql"
	"fmt"

	"social-network/pkg/tools"
)

type LogEntry struct {
	Level    string // "INFO", "ERROR", etc.
	Message  string
	Metadata map[string]any
}

type LoggerModel struct {
	DB *sql.DB
}

func (lm *LoggerModel) Log(entry LogEntry) error {
	metaJson, err := tools.JsonString(entry.Metadata)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	query := `
		INSERT INTO logs (level, message, metadata)
		VALUES (?, ?, ?)`

	if _, err = lm.DB.Exec(query, entry.Level, entry.Message, string(metaJson)); err != nil {
		return fmt.Errorf("insert log: %w", err)
	}

	return nil
}

// Usage:
// logEntry := LogEntry{
// 	Level:   "ERROR",
// 	Message: "Failed to create user",
// 	Metadata: map[string]any{
// 		"user_id": 42,
// 		"ip":      r.RemoteAddr,
// 		"path":    r.URL.Path,
// 	},
// }
// app.Logger.InsertLog(logEntry)
