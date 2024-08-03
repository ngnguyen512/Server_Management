package main

import (

	//_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"server-management/internal/handler"
	"server-management/internal/health_event"
	"server-management/internal/server"
	"server-management/internal/user"
	"strconv"
	"time"

	// "server-management/pkg/repositories"
	"server-management/pkg/encryptoha"
	"server-management/pkg/jwtha"
	"server-management/pkg/postgresha"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"server-management/pkg/loggerha"
)

func main() {
	e := echo.New()

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	//convert port to integer
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	config := &postgresha.Config{
		Host:    os.Getenv("DB_HOST"),
		Port:    port,
		User:    os.Getenv("DB_USER"),
		DBName:  os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	}

	cw := postgresha.NewClientWrapper(config)
	cw.Automigrate(&user.User{}, &server.Server{}, &health_event.HealthEvent{})
	logger := loggerha.NewLogger()

	// Add the custom logger middleware
	e.Use(loggerha.LoggerMiddleware(logger))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// userRepo := repositories.NewUserRepository(db)
	userRepo := postgresha.NewRepository[user.User](cw)

	userHandler := handler.NewUserHandler(userRepo)

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
	config2 := encryptoha.Argon2Config{
		Salt:      []byte("somesalt"),
		Time:      1,
		Memory:    64 * 1024,
		Threads:   4,
		KeyLength: 32,
	}

	encryptor := encryptoha.NewArgon2Encryptor(config2)

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
