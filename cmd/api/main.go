package main

import (
	"context"
	"djiroutine-clean-architecture/internal/http/routes"
	_authUsecase "djiroutine-clean-architecture/internal/modules/auth/usercase"
	_userRepository "djiroutine-clean-architecture/internal/modules/user/repository"
	_userUsecase "djiroutine-clean-architecture/internal/modules/user/usercase"
	"djiroutine-clean-architecture/pkg/config"
	"djiroutine-clean-architecture/pkg/logger"
	"djiroutine-clean-architecture/pkg/sso"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize logger
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	l := logger.L

	// Load environment variables or configuration
	clientID := os.Getenv("OAUTH_CLIENT_ID")
	clientSecret := os.Getenv("OAUTH_CLIENT_SECRET")
	redirectURI := os.Getenv("OAUTH_REDIRECT_URI")
	environment := sso.Environment(os.Getenv("OAUTH_ENVIRONMENT"))
	redisURL := os.Getenv("REDIS_URL")

	timeout, _ := strconv.Atoi(os.Getenv("APP_TIMEOUT"))
	timeoutContext := time.Duration(timeout) * time.Second

	if environment == "" {
		environment = sso.Sandbox // Default to sandbox
	}

	// Initialize OAuth client
	oauthClient, err := sso.NewOAuth2Client(
		clientID,
		clientSecret,
		redirectURI,
		environment,
		redisURL,
	)
	if err != nil {
		log.Fatalf("Failed to initialize OAuth client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize database connection
	pg_port, _ := strconv.Atoi(os.Getenv("DB_PG_PORT"))
	mainDbConfig := config.DBConfig{
		Host:         os.Getenv("DB_PG_HOST"),
		Port:         pg_port,
		User:         os.Getenv("DB_PG_USER"),
		Password:     os.Getenv("DB_PG_PASS"),
		DatabaseName: os.Getenv("DB_PG_DB"),
		MaxConns:     10,
		MinConns:     2,
		IdleTimeout:  5 * time.Minute,
		SSLMode:      "disable",
	}

	mainDbService, err := config.NewDBService(ctx, mainDbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Echo
	e := echo.New()

	// Add standard middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize use cases
	authUseCase := _authUsecase.NewAuthUseCase(oauthClient)

	userRepo := _userRepository.NewUserRepository(mainDbService, l)
	userUsecase := _userUsecase.NewUserUsecase(userRepo, timeoutContext, l)

	useCases := map[string]interface{}{
		"auth": authUseCase,
		"user": userUsecase,
	}

	// Setup routes
	routes.SetupRoutes(e, useCases)

	// Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
