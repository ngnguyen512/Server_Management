package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"server-management/models"
)

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open("postgres", "host=localhost port=5433 user=postgres dbname=servermanagement sslmode=disable password=yourpassword")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Auto migrate structs
	db.AutoMigrate(&models.User{}, &models.Server{}, &models.HealthEvent{})

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
