package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"server-management/internal/handler"
	"server-management/internal/health_event"
	"server-management/internal/server"
	"server-management/internal/user"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5433 user=postgres dbname=servermanagement sslmode=disable password=yourpassword")
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	defer db.Close()

	// Auto migrate structs
	db.AutoMigrate(&user.User{}, &server.Server{}, &health_event.HealthEvent{})

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handler.HelloWorld)

	e.Logger.Fatal(e.Start(":8080"))
}
