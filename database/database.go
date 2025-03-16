package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbConfig struct {
	db  *gorm.DB
	dsn string
}

func (conf *dbConfig) Start() (*gorm.DB, error) {
	// Initialize the database connection.
	db, err := gorm.Open(postgres.Open(conf.dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	fmt.Print("Database connection established successfully..")

	// Automigrate models here.

	return db, nil
}
