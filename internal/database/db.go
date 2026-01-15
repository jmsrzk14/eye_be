package database

import (
	"log"
	"os"
	"time"

	"TA072025/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the database connection
var DB *gorm.DB

func Initialize() {
	var err error

	// 🔑 Ambil DATABASE_URL dari Railway
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL tidak ditemukan")
	}

	// Logger GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)

	// Connect ke PostgreSQL
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Error getting underlying SQL DB: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to PostgreSQL database successfully")

	// Auto migration
	err = DB.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		log.Fatalf("Error auto-migrating database: %v", err)
	}

	log.Println("Database schema migrated successfully")
}

func Close() {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

func GetDB() *gorm.DB {
	return DB
}
