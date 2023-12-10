package server

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"

	"log"
)

func InitDatabase(config *viper.Viper) *sql.DB {
	connectionString := config.GetString("DATABASE_URL")
	maxIdleConnections := config.GetInt("database.max_idle_connections")
	maxOpenConnections := config.GetInt("database.max_open_connections")
	connectionMaxLifetime := config.GetDuration("database.connection_max_lifetime")
	driverName := config.GetString("database.driver_name")

	if connectionString == "" && config.GetString("http.release_mode") == "false" {
		log.Printf("Database connection string is missing, using DATABASE_URL environment variable")
		connectionString = config.GetString("database.url")
	}

	dbHandler, err := sql.Open(driverName, connectionString)
	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}

	dbHandler.SetMaxIdleConns(maxIdleConnections)
	dbHandler.SetMaxOpenConns(maxOpenConnections)
	dbHandler.SetConnMaxLifetime(connectionMaxLifetime)

	err = dbHandler.Ping()
	if err != nil {
		dbHandler.Close()
		log.Fatalf("Error while validating database: %v", err)
	}

	return dbHandler
}
