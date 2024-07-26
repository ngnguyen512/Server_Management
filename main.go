package main

import (

	//_ "github.com/jinzhu/gorm/dialects/postgres"
	"server-management/internal/handler"
	"server-management/internal/health_event"
	"server-management/internal/server"
	"server-management/internal/user"

	// "server-management/pkg/repositories"
	"server-management/pkg/postgresha"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	cw := postgresha.NewClientWrapper()
	cw.Automigrate(&user.User{}, &server.Server{}, &health_event.HealthEvent{})

	// userRepo := repositories.NewUserRepository(db)
	userRepo := postgresha.NewRepository[user.User](cw)

	userHandler := handler.NewUserHandler(userRepo)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", handler.HelloWorld)

	e.POST("/users", userHandler.CreateUser)
	e.GET("/users/:id", userHandler.GetUserById)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)

	serverRepo := postgresha.NewRepository[server.Server](cw)
	serverHandler := handler.NewServerHandler(serverRepo)

	e.POST("/servers", serverHandler.CreateServer)
	e.GET("/servers/:id", serverHandler.GetServerById)
	e.PUT("/servers/:id", serverHandler.UpdateServer)
	e.DELETE("/servers/:id", serverHandler.DeleteOneById)

	heathRepo := postgresha.NewRepository[health_event.HealthEvent](cw)
	healthHandler := handler.NewHealthHandler(heathRepo)

	e.POST("/health-events", healthHandler.CreateHealthEvent)
	e.GET("/health-events/:id", healthHandler.GetHealthEventById)
	e.PUT("/health-events/:id", healthHandler.UpdateHealthEvent)
	e.DELETE("/health-events/:id", healthHandler.DeleteHealthEvent)

	e.Logger.Fatal(e.Start(":8080"))
}
