package model

import (
	"strings"

	"github.com/jinzhu/gorm"

	// Postgres Driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// PSQL is a database object that holds connection to PostgreSQL database.
var PSQL *gorm.DB

// InitPSQLConnection establishes connection to database.
func InitPSQLConnection() (*gorm.DB, error) {
	if PSQL != nil {
		return PSQL, nil
	}

	credentials := []string{
		"host=localhost",
		"port=5432",
		"user=crawler",
		"password=crawler",
		"dbname=moviecrawler",
		"sslmode=disable",
	}

	var err error
	PSQL, err = gorm.Open("postgres", strings.Join(credentials, " "))
	if err != nil {
		return nil, err
	}

	PSQL.LogMode(false)

	return PSQL, nil
}
