package database

import (
	"context"
	"time"

	"github.com/Dom-HTG/attendance-management-system/entities"
	"github.com/Dom-HTG/attendance-management-system/pkg/logger"
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
	if err != nil {
		return nil, err
	}

	// DO NOT close postgresDB here: closing the underlying sql.DB will
	// invalidate the returned *gorm.DB and cause "sql: database is closed"
	// errors when the application later tries to use the DB connection.

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
		logger.Errorf("database ping failed: %v", e)
		return nil, e
	}

	logger.Info("Database connection established successfully")

	// Migrate models and report outcome
	if err := db.AutoMigrate(
		&entities.Student{},
		&entities.Lecturer{},
		&entities.Event{},
		&entities.UserAttendance{},
		&entities.Admin{},
		&entities.AuditLog{},
		&entities.SystemSettings{},
	); err != nil {
		logger.Errorf("AutoMigrate failed: %v", err)
		return nil, err
	}

	logger.Info("Database migrations applied successfully")

	return db, nil
}
