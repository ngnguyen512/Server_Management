package postgresha

import (
	"database/sql"
	"log"
	"time"
)

type ClientWrapper struct {
	connectionString string
	db               *sql.DB
}

func NewClientWrapper() *ClientWrapper {
	connectionString := "host=localhost port=5433 user=postgres dbname=servermanagement sslmode=disable"
	return &ClientWrapper{connectionString: connectionString}
}

func (cw *ClientWrapper) connect() error {
	var err error
	cw.db, err = sql.Open("postgres", cw.connectionString)
	if err != nil {
		return err
	}
	return cw.db.Ping()
}

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

func (cw *ClientWrapper) Db() *sql.DB {
	if cw.db == nil {
		cw.doConnect()
	}
	return cw.db
}
