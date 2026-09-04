package model

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	db, err := gorm.Open(postgres.Open(databaseURL()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("open postgres database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get postgres sql database: %w", err)
	}

	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.PingContext(context.Background()); err != nil {
		_ = sqlDB.Close()
		return fmt.Errorf("ping postgres: %w", err)
	}

	DB = db
	return nil
}

func CloseDB() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func databaseURL() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}

	dsn := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(envOrDefault("DB_USER", "postgres"), envOrDefault("DB_PASSWORD", "postgres")),
		Host:   net.JoinHostPort(envOrDefault("DB_HOST", "localhost"), envOrDefault("DB_PORT", "5432")),
		Path:   envOrDefault("DB_NAME", "taskDB"),
	}
	query := dsn.Query()
	query.Set("sslmode", envOrDefault("DB_SSLMODE", "disable"))
	dsn.RawQuery = query.Encode()
	return dsn.String()
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
