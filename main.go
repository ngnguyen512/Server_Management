package main

import (

	//_ "github.com/jinzhu/gorm/dialects/postgres"
	"server-management/internal/handler"
	"server-management/internal/health_event"
	"server-management/internal/server"
	"server-management/internal/user"
	"time"

	// "server-management/pkg/repositories"
	"server-management/pkg/encryptoha"
	"server-management/pkg/jwtha"
	"server-management/pkg/postgresha"

	"github.com/golang-jwt/jwt"
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

	g2 := e.Group("/servers")

	serverRepo := postgresha.NewRepository[server.Server](cw)
	serverHandler := handler.NewServerHandler(serverRepo)

	g2.POST("/servers", serverHandler.CreateServer)
	g2.GET("/servers/:id", serverHandler.GetServerById)
	g2.PUT("/servers/:id", serverHandler.UpdateServer)
	g2.DELETE("/servers/:id", serverHandler.DeleteOneById)

	heathRepo := postgresha.NewRepository[health_event.HealthEvent](cw)
	healthHandler := handler.NewHealthHandler(heathRepo)

	g1 := e.Group("/health-events")
	g1.POST("/health-events", healthHandler.CreateHealthEvent)
	g1.GET("/health-events/:id", healthHandler.GetHealthEventById)
	g1.PUT("/health-events/:id", healthHandler.UpdateHealthEvent)
	g1.DELETE("/health-events/:id", healthHandler.DeleteHealthEvent)

	g := e.Group("/users")
	g.POST("/", userHandler.CreateUser)
	g.GET("/:id", userHandler.GetUserById)
	g.PUT("/:id", userHandler.UpdateUser)
	g.DELETE("/:id", userHandler.DeleteUser)
	config := encryptoha.Argon2Config{
		Salt:      []byte("somesalt"),
		Time:      1,
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: 32,
	}

	encryptor := encryptoha.NewArgon2Encryptor(config)

	jwtConfig := jwtha.JwtConfig{
		SecretKey:     []byte("supersecretkey"),
		SigningMethod: jwt.SigningMethodHS256,
		Expiration:    time.Hour,
	}
	jwtTokenizer := jwtha.NewJwtTokenizer(jwtConfig)

	authHandler := handler.NewAuthHandler(encryptor, *jwtTokenizer, userRepo)
	e.POST("/auth/signup", authHandler.SignUp)
	e.POST("/auth/login", authHandler.Login)

	e.Logger.Fatal(e.Start(":8080"))

}
