package database

import (
	"database/sql"
	"fmt"
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
	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("Failed to ping database: %v", err)
	}
	return db, nil

}
