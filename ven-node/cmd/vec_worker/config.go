package main

import (
	"os"
)

type Config struct {
	MySQLHost     string
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string
	MySQLPort     string
}

func loadConfig() (*Config, error) {
	cfg := &Config{
		MySQLHost:     os.Getenv("MYSQL_HOST"),
		MySQLUser:     os.Getenv("MYSQL_USER"),
		MySQLPassword: os.Getenv("MYSQL_PASSWORD"),
		MySQLDatabase: os.Getenv("MYSQL_DBNAME"),
		MySQLPort:     os.Getenv("MYSQL_PORT"),
	}

	// TODO add validation logic here to check if required environment variables are set.

	return cfg, nil
}
