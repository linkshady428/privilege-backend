package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/privilege/backend/internal/config"
)

type Claims struct {
	UserID string `json:"user_id"`
	Tier   string `json:"tier"` // "free" or "privilege"
	jwt.RegisteredClaims
}

// ContextKeyUserID is the echo context key for the authenticated user.
const ContextKeyUserID = "user_id"
const ContextKeyTier = "tier"

// JWT returns an Echo middleware that validates Bearer JWT tokens.
// Set SKIP_AUTH=true to bypass validation during local development.
func JWT(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.SkipAuth {
				c.Set(ContextKeyUserID, "dev-user-id")
				c.Set(ContextKeyTier, "privilege")
				return next(c)
			}

			raw := c.Request().Header.Get(echo.HeaderAuthorization)
			if !strings.HasPrefix(raw, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
			}
			tokenStr := strings.TrimPrefix(raw, "Bearer ")

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(cfg.JWTSecret), nil
			})
			if err != nil || !token.Valid {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			c.Set(ContextKeyUserID, claims.UserID)
			c.Set(ContextKeyTier, claims.Tier)
			return next(c)
		}
	}
}

// RequireTier returns a middleware that gates access to a specific tier.
func RequireTier(tier string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Get(ContextKeyTier) != tier {
				return echo.NewHTTPError(http.StatusForbidden, "requires "+tier+" tier")
			}
			return next(c)
		}
	}
}
