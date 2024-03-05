package database

import (
	"apk-sekolah/config"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func configureDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(500)
	sqlDB.SetConnMaxLifetime(time.Second * 5)

	return db, nil
}

func InitPostgreSQL(appConfig *config.AppConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable timezone=Asia/Shanghai",
		appConfig.DBHost, appConfig.DBUsername, appConfig.DBPassword, appConfig.DBName, appConfig.DBPort)

	log.Printf("Using DSN: %s", dsn)

	DB, err := configureDB(dsn)
	if err != nil {
		log.Printf("Failed to initialize database connection: %s", err)
		return nil, err
	}

	log.Println("Connected to the database")

	return DB, nil
}
