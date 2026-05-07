package handler

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	mw "github.com/privilege/backend/internal/middleware"
)

type FeedHandler struct {
	db *pgxpool.Pool
}

func NewFeedHandler(pool *pgxpool.Pool) *FeedHandler { return &FeedHandler{db: pool} }

type FeedCard struct {
	UserID string   `json:"user_id"`
	Name   string   `json:"name"`
	Age    int      `json:"age"`
	Bio    string   `json:"bio"`
	Photos []string `json:"photos"`
}

type InvitationCard struct {
	InvitationID string    `json:"invitation_id"`
	Sender       *FeedCard `json:"sender,omitempty"`
	Sealed       bool      `json:"sealed"`
}

// GET /api/v1/feed
// Privilege: card stack of nearby Free users (last active, unseen only).
// Free: sealed invitation stack.
func (h *FeedHandler) GetFeed(c echo.Context) error {
	userID := c.Get(mw.ContextKeyUserID).(string)
	tier := c.Get(mw.ContextKeyTier).(string)
	ctx := c.Request().Context()

	if tier == "privilege" {
		rows, err := h.db.Query(ctx, `
			SELECT id, name, birthdate, COALESCE(bio, '')
			FROM users
			WHERE tier = 'free'
			  AND deleted_at IS NULL
			  AND id NOT IN (
			    SELECT free_user_id FROM passes WHERE privilege_user_id = $1
			  )
			ORDER BY last_active DESC
			LIMIT 20`, userID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "server error")
		}
		defer rows.Close()

		cards := []FeedCard{}
		for rows.Next() {
			var card FeedCard
			var birthdate time.Time
			if err := rows.Scan(&card.UserID, &card.Name, &birthdate, &card.Bio); err != nil {
				continue
			}
			card.Age = int(time.Since(birthdate).Hours() / (24 * 365.25))
			card.Photos = []string{}
			cards = append(cards, card)
		}
		return c.JSON(http.StatusOK, echo.Map{"cards": cards})
	}

	// Free user — pending invitations (sealed: sender identity hidden until accepted)
	rows, err := h.db.Query(ctx, `
		SELECT id FROM invitations
		WHERE recipient_id = $1 AND status = 'pending'
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	}
	defer rows.Close()

	invitations := []InvitationCard{}
	for rows.Next() {
		var card InvitationCard
		rows.Scan(&card.InvitationID)
		card.Sealed = true
		invitations = append(invitations, card)
	}
	return c.JSON(http.StatusOK, echo.Map{"invitations": invitations})
}

// POST /api/v1/feed/pass/:userID  (Privilege tier only)
func (h *FeedHandler) PassUser(c echo.Context) error {
	freeUserID := c.Param("userID")
	privilegeUserID := c.Get(mw.ContextKeyUserID).(string)

	if _, err := h.db.Exec(c.Request().Context(),
		`INSERT INTO passes (privilege_user_id, free_user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		privilegeUserID, freeUserID,
	); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "passed"})
}
