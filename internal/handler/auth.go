package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler { return &AuthHandler{} }

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"` // ISO 8601 date, e.g. "1995-06-15"
	Sex       string `json:"sex"`
	AgreeToS  bool   `json:"agree_tos"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Tier         string `json:"tier"`
}

// POST /api/v1/auth/register
func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// TODO: validate age gate (must be 18+), hash password, persist user
	return c.JSON(http.StatusCreated, AuthResponse{
		AccessToken:  "stub-access-token",
		RefreshToken: "stub-refresh-token",
		Tier:         "free",
	})
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// TODO: verify credentials, issue JWT pair
	return c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  "stub-access-token",
		RefreshToken: "stub-refresh-token",
		Tier:         "free",
	})
}

// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// TODO: validate refresh token, rotate it, issue new access token
	return c.JSON(http.StatusOK, AuthResponse{
		AccessToken:  "stub-access-token",
		RefreshToken: "stub-new-refresh-token",
		Tier:         "free",
	})
}

// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c echo.Context) error {
	// TODO: revoke refresh token
	return c.JSON(http.StatusOK, echo.Map{"message": "logged out"})
}
