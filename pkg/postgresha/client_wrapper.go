package postgresha

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host    string
	Port    int
	User    string
	DBName  string
	SSLMode string
}

type ClientWrapper struct {
	config *Config
	db     *gorm.DB
}

// NewClientWrapper creates a new instance of ClientWrapper
func NewClientWrapper(config *Config) *ClientWrapper {
	return &ClientWrapper{config: config}
}

func (cw *ClientWrapper) connectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s",
		cw.config.Host, cw.config.Port, cw.config.User, cw.config.DBName, cw.config.SSLMode)
}
func (cw *ClientWrapper) Automigrate(d ...interface{}) error {
	err := cw.Db().AutoMigrate(d...)
	return err
}

// connect uses GORM to open a database connection
func (cw *ClientWrapper) connect() error {
	var err error
	cw.db, err = gorm.Open(postgres.Open(cw.connectionString()), &gorm.Config{})
	if err != nil {
		return err
	}

	// Ping the database to check connection
	sqlDB, err := cw.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// doConnect manages retries for establishing the database connection
func (cw *ClientWrapper) doConnect() {
	for {
		if err := cw.connect(); err != nil {
			log.Printf("Failed to connect to database: %v. Retrying...", err)
			time.Sleep(5 * time.Second) // Retry after 5 seconds
		} else {
			break
		}
	}
}

// Db returns the GORM database object
func (cw *ClientWrapper) Db() *gorm.DB {
	if cw.db == nil {
		cw.doConnect()
	}
	return cw.db
}
