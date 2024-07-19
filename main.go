package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	//_ "github.com/jinzhu/gorm/dialects/postgres"
	"server-management/internal/handler"
	"server-management/internal/health_event"
	"server-management/internal/server"
	"server-management/internal/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var db *gorm.DB

func main() {
	var err error
	dsn := "host=localhost user=postgres password=yourpassword dbname=servermanagement port=5433 sslmode=disable TimeZone=Asia/Shanghai"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get generic database object: " + err.Error())
	}
	defer sqlDB.Close()

	// Auto migrate structs
	err = db.AutoMigrate(&user.User{}, &server.Server{}, &health_event.HealthEvent{})

	if err != nil {
		panic("failed to auto-migrate: " + err.Error())
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handler.HelloWorld)

	e.Logger.Fatal(e.Start(":8080"))
}
