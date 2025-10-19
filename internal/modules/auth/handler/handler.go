package handler

import (
	"djiroutine-clean-architecture/internal/modules/auth"
	"djiroutine-clean-architecture/pkg/errors"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authUseCase auth.UseCase
}

func NewAuthHandler(authUseCase auth.UseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

// Login mengarahkan pengguna ke halaman login OAuth
func (h *AuthHandler) Login(c echo.Context) error {
	authURL, state, err := h.authUseCase.GetAuthorizationURL(c.Request().Context())
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return c.JSON(appErr.Code, map[string]string{
				"error": appErr.Message,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate authorization URL",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"auth_url": authURL,
		"state":    state,
	})
}

// Callback menangani callback dari OAuth provider
func (h *AuthHandler) Callback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" || state == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing code or state parameter",
		})
	}

	user, token, err := h.authUseCase.ProcessCallback(c.Request().Context(), code, state)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return c.JSON(appErr.Code, map[string]string{
				"error": appErr.Message,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Authentication failed",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

// Logout mengakhiri sesi pengguna
func (h *AuthHandler) Logout(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Authorization header is required",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Authorization header format must be Bearer {token}",
		})
	}

	token := parts[1]

	err := h.authUseCase.Logout(c.Request().Context(), token)
	if err != nil {
		if appErr, ok := err.(*errors.AppError); ok {
			return c.JSON(appErr.Code, map[string]string{
				"error": appErr.Message,
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Logout failed",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Successfully logged out",
	})
}
