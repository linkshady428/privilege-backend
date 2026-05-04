package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	mw "github.com/privilege/backend/internal/middleware"
)

type FeedHandler struct{}

func NewFeedHandler() *FeedHandler { return &FeedHandler{} }

type FeedCard struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Bio    string `json:"bio"`
	Photos []string `json:"photos"`
}

type InvitationCard struct {
	InvitationID string `json:"invitation_id"`
	// Privilege user details revealed on tap
	Sender *FeedCard `json:"sender,omitempty"`
	Sealed bool      `json:"sealed"`
}

// GET /api/v1/feed
// Privilege users get a card stack of nearby Free users (ordered by last active, unseen only).
// Free users get a stack of sealed invitation cards.
func (h *FeedHandler) GetFeed(c echo.Context) error {
	tier := c.Get(mw.ContextKeyTier).(string)

	if tier == "privilege" {
		cards := []FeedCard{
			{UserID: "stub-free-user-1", Name: "Alex", Age: 24, Bio: "Coffee lover", Photos: []string{}},
		}
		return c.JSON(http.StatusOK, echo.Map{"cards": cards})
	}

	// Free user — sealed invitation stack
	invitations := []InvitationCard{
		{InvitationID: "stub-inv-1", Sealed: true},
	}
	return c.JSON(http.StatusOK, echo.Map{"invitations": invitations})
}

// POST /api/v1/feed/pass/:userID  (Privilege tier only)
// Records a permanent left-swipe so the Free user never reappears.
func (h *FeedHandler) PassUser(c echo.Context) error {
	userID := c.Param("userID")
	_ = userID
	// TODO: append to privilege user's pass log (append-only, no FK delete)
	return c.JSON(http.StatusOK, echo.Map{"message": "passed"})
}
