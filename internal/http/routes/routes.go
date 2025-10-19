package routes

import (
	"djiroutine-go-clean-architecture/internal/http/middleware"
	"djiroutine-go-clean-architecture/internal/modules/auth"
	authHandler "djiroutine-go-clean-architecture/internal/modules/auth/handler"
	"djiroutine-go-clean-architecture/internal/modules/user"
	userHandler "djiroutine-go-clean-architecture/internal/modules/user/handler"
	"djiroutine-go-clean-architecture/pkg/logger"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, useCases map[string]interface{}) {
	authUseCase, ok := useCases["auth"].(auth.UseCase)
	if !ok {
		panic("Invalid auth use case provided")
	}

	authH := authHandler.NewAuthHandler(authUseCase)

	oauthMiddleware := middleware.NewOAuthMiddleware(authUseCase)

	authGroup := e.Group("/auth")
	authGroup.GET("/login", authH.Login)
	authGroup.GET("/callback", authH.Callback)
	authGroup.POST("/logout", authH.Logout)

	apiGroup := e.Group("/api")
	apiGroup.Use(oauthMiddleware.Authenticate)

	apiGroup.GET("/hello", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	setupUsersRoutes(apiGroup, useCases)
}

func setupUsersRoutes(g *echo.Group, useCases map[string]interface{}) {
	userUseCase, ok := useCases["user"].(user.UseCase)
	if !ok {
		panic("Invalid user use case provided")
	}

	userH := userHandler.NewUserHandler(logger.L, userUseCase)
	g.GET("/users", userH.ListUsers)
}
