package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Notification struct {
	ID        string
	Type      string
	Message   string
	Sent      bool
	CreatedAt time.Time
}

type DB struct {
	conn *sql.DB
}

var (
	instance *DB
	once     sync.Once
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetInstance() (*DB, error) {
	var err error
	once.Do(func() {
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		dbname := getEnv("DB_NAME", "testworkflow")
		user := getEnv("DB_USER", "postgres")
		password := getEnv("DB_PASSWORD", "postgres")

		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)

		var conn *sql.DB
		conn, err = sql.Open("postgres", connStr)
		if err != nil {
			return
		}

		err = conn.Ping()
		if err != nil {
			return
		}

		instance = &DB{conn: conn}
		err = instance.createSchema()
	})
	return instance, err
}

func (db *DB) createSchema() error {
	query := `
		CREATE TABLE IF NOT EXISTS notifications (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			type VARCHAR NOT NULL,
			message TEXT NOT NULL,
			sent BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := db.conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create notifications table: %w", err)
	}
	log.Println("Database schema initialized")
	return nil
}

func (db *DB) SaveNotification(notifType, message string) (string, error) {
	id := uuid.New().String()
	query := `INSERT INTO notifications (id, type, message, sent, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.conn.Exec(query, id, notifType, message, false, time.Now())
	if err != nil {
		return "", fmt.Errorf("failed to save notification: %w", err)
	}
	return id, nil
}

func (db *DB) GetNotifications(limit, offset int32) ([]Notification, int32, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	var total int32
	err := db.conn.QueryRow("SELECT COUNT(*) FROM notifications").Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	query := `SELECT id, type, message, sent, created_at FROM notifications ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := db.conn.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var n Notification
		if err := rows.Scan(&n.ID, &n.Type, &n.Message, &n.Sent, &n.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}

	return notifications, total, nil
}

func (db *DB) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}
