package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func InitializePostgresDb() (*sql.DB, error) {
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbDatabase := getEnvOrDefault("DB_DATABASE", "postgres")
	dbUsername := getEnvOrDefault("DB_USERNAME", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbSslMode := getEnvOrDefault("DB_SSL_MODE", "require")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", dbUsername, dbPassword, dbHost, dbPort, dbDatabase, dbSslMode)
	log.Printf("Connecting to database")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	log.Printf("Pinging database")
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to database")
	return db, nil
}

func CloseDb(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Printf("Error closing dynamo %v", err)
	}
}

func getEnvOrDefault(key string, def string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}

	return def
}
