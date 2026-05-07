package handler

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	mw "github.com/privilege/backend/internal/middleware"
)

type AuthHandler struct {
	db        *pgxpool.Pool
	jwtSecret []byte
}

func NewAuthHandler(pool *pgxpool.Pool, jwtSecret string) *AuthHandler {
	return &AuthHandler{db: pool, jwtSecret: []byte(jwtSecret)}
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Birthdate string `json:"birthdate"` // YYYY-MM-DD
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
	if req.Email == "" || req.Password == "" || req.Name == "" || req.Birthdate == "" || req.Sex == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}
	if !req.AgreeToS {
		return echo.NewHTTPError(http.StatusBadRequest, "must agree to terms of service")
	}

	dob, err := time.Parse("2006-01-02", req.Birthdate)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid birthdate, use YYYY-MM-DD")
	}
	if time.Since(dob) < 18*365*24*time.Hour {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "must be 18 or older")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	}

	var userID, tier string
	err = h.db.QueryRow(c.Request().Context(),
		`INSERT INTO users (email, password_hash, name, birthdate, sex, agree_tos)
		 VALUES ($1, $2, $3, $4, $5::sex_type, $6)
		 RETURNING id, tier::text`,
		req.Email, string(hash), req.Name, dob, req.Sex, req.AgreeToS,
	).Scan(&userID, &tier)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return echo.NewHTTPError(http.StatusConflict, "email already registered")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	}

	return h.issueTokens(c, userID, tier, http.StatusCreated)
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var userID, tier, passwordHash string
	err := h.db.QueryRow(c.Request().Context(),
		`SELECT id, tier::text, COALESCE(password_hash, '')
		 FROM users WHERE email = $1 AND deleted_at IS NULL`,
		strings.ToLower(req.Email),
	).Scan(&userID, &tier, &passwordHash)
	if err == pgx.ErrNoRows || passwordHash == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	return h.issueTokens(c, userID, tier, http.StatusOK)
}

// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	raw, err := hex.DecodeString(req.RefreshToken)
	if err != nil || len(raw) != 32 {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid refresh token")
	}
	sum := sha256.Sum256(raw)
	tokenHash := hex.EncodeToString(sum[:])

	var userID string
	err = h.db.QueryRow(c.Request().Context(),
		`UPDATE refresh_tokens SET revoked = true
		 WHERE token_hash = $1 AND revoked = false AND expires_at > NOW()
		 RETURNING user_id`,
		tokenHash,
	).Scan(&userID)
	if err == pgx.ErrNoRows {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired refresh token")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	}

	var tier string
	if err := h.db.QueryRow(c.Request().Context(),
		`SELECT tier::text FROM users WHERE id = $1`, userID,
	).Scan(&tier); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	}

	return h.issueTokens(c, userID, tier, http.StatusOK)
}

// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c echo.Context) error {
	var req RefreshRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if raw, err := hex.DecodeString(req.RefreshToken); err == nil && len(raw) == 32 {
		sum := sha256.Sum256(raw)
		h.db.Exec(c.Request().Context(),
			`UPDATE refresh_tokens SET revoked = true WHERE token_hash = $1`,
			hex.EncodeToString(sum[:]),
		)
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "logged out"})
}

func (h *AuthHandler) issueTokens(c echo.Context, userID, tier string, status int) error {
	accessToken, err := h.signJWT(userID, tier)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	}
	refreshToken := hex.EncodeToString(raw)
	sum := sha256.Sum256(raw)

	if _, err := h.db.Exec(c.Request().Context(),
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, hex.EncodeToString(sum[:]), time.Now().Add(7*24*time.Hour),
	); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	}

	return c.JSON(status, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Tier:         tier,
	})
}

func (h *AuthHandler) signJWT(userID, tier string) (string, error) {
	claims := mw.Claims{
		UserID: userID,
		Tier:   tier,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(h.jwtSecret)
}
