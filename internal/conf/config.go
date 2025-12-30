package conf

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Conf struct {
	log *logrus.Logger
	db  *gorm.DB
}

func NewConf(log *logrus.Logger, fileName ...string) *Conf {
	if err := godotenv.Load(fileName...); err != nil {
		log.Error("No .env file found\n-> ", err)
		os.Exit(1)
	}

	// Валидация обязательных переменных окружения
	if err := validateRequiredEnvVars(); err != nil {
		log.Fatal(err)
	}

	return &Conf{log: log, db: nil}
}

func (c *Conf) GetConValue(key string) string {
	// Retrieve specific configuration value by key
	return os.Getenv(key)
}

// GetJWTSecret возвращает JWT секрет из окружения
func (c *Conf) GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}

// validateRequiredEnvVars проверяет наличие всех обязательных переменных окружения
func validateRequiredEnvVars() error {
	requiredVars := []string{
		"APP_HOST",
		"APP_PORT",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_NAME",
		"DB_PASSWORD",
		"JWT_SECRET",
	}

	var missing []string
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			missing = append(missing, v)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missing)
	}

	// Дополнительная проверка JWT_SECRET на минимальную длину
	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters long, got %d", len(jwtSecret))
	}

	log.Info("All required environment variables are set")
	return nil
}
func (c *Conf) DBConnect() *gorm.DB {
	host := c.GetConValue("DB_HOST")
	port := c.GetConValue("DB_PORT")
	user := c.GetConValue("DB_USER")
	dbname := c.GetConValue("DB_NAME")
	password := c.GetConValue("DB_PASSWORD")
	sslmode := c.GetConValue("DB_SSLMODE")

	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		c.log.Fatal("Failed to connect to database: ", err)
		os.Exit(1)
	}
	c.log.Info("Database connection established")

	// Set up connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		c.log.Fatal("Failed to get database instance: ", err)
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Ping the database to verify connection
	if err := sqlDB.Ping(); err != nil {
		c.log.Fatal("Failed to ping database: ", err)
		os.Exit(1)
	}

	return db
}
func (c *Conf) ConnectToDb(db_name string) *gorm.DB {
	host := c.GetConValue("DB_HOST")
	port := c.GetConValue("DB_PORT")
	user := c.GetConValue("DB_USER")
	dbname := db_name
	password := c.GetConValue("DB_PASSWORD")
	sslmode := c.GetConValue("DB_SSLMODE")

	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		c.log.Fatal("Failed to connect to database: ", err)
		os.Exit(1)
	}
	c.log.Info("Database connection established")

	// Set up connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		c.log.Fatal("Failed to get database instance: ", err)
		os.Exit(1)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	// 	// Ping the database to verify connection
	if err := sqlDB.Ping(); err != nil {
		c.log.Fatal("Failed to ping database: ", err)
		os.Exit(1)
	}

	return db
}

// git remote add origin https://github.com/rusgainew/tunduct-project-system.git
// git branch -M dev
// git push -u origin dev
