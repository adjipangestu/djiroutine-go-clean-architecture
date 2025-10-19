package middleware

import (
	"djiroutine-go-clean-architecture/internal/modules/auth"
	"djiroutine-go-clean-architecture/pkg/errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type OAuthMiddleware struct {
	AuthUseCase auth.UseCase
}

func NewOAuthMiddleware(authUseCase auth.UseCase) *OAuthMiddleware {
	return &OAuthMiddleware{
		AuthUseCase: authUseCase,
	}
}

func (m *OAuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Authorization header is required",
			})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Authorization header format must be Bearer {token}",
			})
		}

		token := parts[1]

		// Verifikasi token melalui use case
		user, err := m.AuthUseCase.ValidateToken(c.Request().Context(), token)
		fmt.Println(user)
		if err != nil {
			if appErr, ok := err.(*errors.AppError); ok {
				return c.JSON(appErr.Code, map[string]string{
					"error": appErr.Message,
				})
			}
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "Invalid or expired token",
			})
		}

		// Tambahkan user ke context
		c.Set("user", user)

		return next(c)
	}
}
