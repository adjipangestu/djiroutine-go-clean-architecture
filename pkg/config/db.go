package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	User         string        `json:"user"`
	Password     string        `json:"password"`
	DatabaseName string        `json:"database_name"`
	MaxConns     int           `json:"max_conns"`
	MinConns     int           `json:"min_conns"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
	SSLMode      string        `json:"ssl_mode"`
}

type DBService interface {
	GetConnection() *gorm.DB
	Close()
}

type DBServiceImpl struct {
	db *gorm.DB
}

func NewDBService(ctx context.Context, config DBConfig) (DBService, error) {
	sslMode := config.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DatabaseName, sslMode,
	)

	// Buat context dengan timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("❌ gagal koneksi ke database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("❌ gagal mengonfigurasi database: %w", err)
	}

	// Konfigurasi koneksi
	sqlDB.SetMaxOpenConns(config.MaxConns)
	sqlDB.SetMaxIdleConns(config.MinConns)
	sqlDB.SetConnMaxIdleTime(config.IdleTimeout)

	// Cek koneksi dengan Ping()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("❌ gagal ping database: %w", err)
	}

	log.Println("✅ Database terkoneksi dengan sukses!")
	return &DBServiceImpl{db: db}, nil
}

func (db *DBServiceImpl) GetConnection() *gorm.DB {
	return db.db
}

func (db *DBServiceImpl) Close() {
	sqlDB, err := db.db.DB()
	if err != nil {
		log.Println("⚠️ Gagal menutup koneksi database:", err)
		return
	}
	sqlDB.Close()
	log.Println("🔴 Koneksi database ditutup")
}
