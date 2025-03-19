package database

import (
	"context"
	"fmt"
	"time"

	"github.com/Dom-HTG/attendance-management-system/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbConfig struct {
	DSN           string
	MaxOpenConns  int
	MaxIdleConns  int
	MaxIdleTimout string
}

func (conf *DbConfig) Start() (*gorm.DB, error) {
	// Initialize the database connection.
	db, err := gorm.Open(postgres.Open(conf.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Set pooling configuration.
	postgresDB, err := db.DB()
	defer postgresDB.Close()

	postgresDB.SetMaxOpenConns(conf.MaxOpenConns)
	postgresDB.SetMaxIdleConns(conf.MaxIdleConns)

	lifetime, er := time.ParseDuration(conf.MaxIdleTimout)
	if er != nil {
		return nil, er
	}
	postgresDB.SetConnMaxLifetime(lifetime)

	// ping database.
	ctx, canc := context.WithTimeout(context.Background(), time.Second*3) // timeout in 3s.
	defer canc()

	e := postgresDB.PingContext(ctx)
	if e != nil {
		return nil, e
	}

	fmt.Print("Database connection established successfully..")

	// Migrate models.
	db.AutoMigrate(
		&entities.Student{},
		&entities.Lecturer{},
		&entities.Event{},
		&entities.Attendance{},
		&entities.UserAttendance{},
	)

	return db, nil
}
