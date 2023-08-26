package linkcxo

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// UserCredential -
type UserCredential struct {
	IsService   bool     `json:"isService"`
	UserID      string   `json:"userId"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	Token       string   `json:"token"`
}

func LoggingMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		res := next(c)

		logrus.WithFields(logrus.Fields{
			"method":     c.Request().Method,
			"path":       c.Path(),
			"status":     c.Response().Status,
			"latency_ns": time.Since(start).Nanoseconds(),
		}).Info("request details")

		return res
	}
}

type IAuthService interface {
	Authenticate(token string) (*UserCredential, error)
}

type AuthConfig struct {
	Skipper     func(c echo.Context) bool
	AuthService IAuthService
}

// AuthMiddleware - Authenticate
func AuthMiddleware(authService IAuthService) echo.MiddlewareFunc {
	return AuthMiddlewareWithConfig(AuthConfig{
		Skipper: func(c echo.Context) bool {
			return false
		},
		AuthService: authService,
	})
}
func AuthMiddlewareWithConfig(config AuthConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if !config.Skipper(c) {
				bearerToken := c.Request().Header.Get(echo.HeaderAuthorization)
				if bearerToken != "" {
					userCred, err := config.AuthService.Authenticate(bearerToken)
					if err != nil {
						logrus.Errorln(err)
						return c.JSON(http.StatusUnauthorized, NewResponseBuilder().BuildError(errors.New(err.Error()), ErrorCode.Common.StatusUnauthorized, http.StatusUnauthorized))
					}
					RequestUtils{}.SetCredential(c, *userCred)
				} else {
					return c.JSON(http.StatusUnauthorized, NewResponseBuilder().BuildError(errors.New("Unauthorized"), ErrorCode.Common.StatusUnauthorized, http.StatusUnauthorized))
				}
			}
			return next(c)
		}
	}
}
