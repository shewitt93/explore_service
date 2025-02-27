package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math"
	"time"
)

type ConfigDatabase struct {
	User     string
	Password string
	Host     string
	Port     string
}

func GenerateDSN(config ConfigDatabase, database string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.User, config.Password, config.Host, config.Port, database)
}

func NewMysqlConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to database: %v", err)
	}

	// Try to ping with retries
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		err = db.Ping()
		if err == nil {
			// Connection successful
			break
		}

		log.Printf("Database ping attempt %d/%d failed: %v", i+1, maxRetries, err)

		if i < maxRetries-1 {
			// Wait before retrying - exponential backoff
			waitTime := time.Duration(math.Pow(2, float64(i))) * time.Second
			log.Printf("Waiting %v before next attempt...", waitTime)
			time.Sleep(waitTime)
		}
	}

	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("Failed to ping database after %d attempts: %v", maxRetries, err)
	}

	log.Println("Successfully connected to the database")
	return db, nil
}
