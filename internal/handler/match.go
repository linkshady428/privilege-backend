package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type MatchHandler struct{}

func NewMatchHandler() *MatchHandler { return &MatchHandler{} }

type Match struct {
	MatchID   string    `json:"match_id"`
	Partner   *FeedCard `json:"partner"`
	CreatedAt string    `json:"created_at"`
}

// GET /api/v1/matches
func (h *MatchHandler) ListMatches(c echo.Context) error {
	// TODO: query matches for current user, join partner profile
	matches := []Match{
		{
			MatchID:   "stub-match-1",
			Partner:   &FeedCard{UserID: "partner-1", Name: "Jordan", Age: 27, Bio: "", Photos: []string{}},
			CreatedAt: "2024-01-01T00:00:00Z",
		},
	}
	return c.JSON(http.StatusOK, echo.Map{"matches": matches})
}
