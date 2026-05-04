package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type InvitationHandler struct{}

func NewInvitationHandler() *InvitationHandler { return &InvitationHandler{} }

type SendInvitationRequest struct {
	RecipientUserID string `json:"recipient_user_id"`
}

// POST /api/v1/invitations  (Privilege tier only — swipe right)
func (h *InvitationHandler) SendInvitation(c echo.Context) error {
	var req SendInvitationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// TODO: create invitation record, send push notification to Free user
	return c.JSON(http.StatusCreated, echo.Map{
		"invitation_id": "stub-inv-id",
		"status":        "pending",
	})
}

// POST /api/v1/invitations/:invitationID/accept  (Free tier only)
func (h *InvitationHandler) AcceptInvitation(c echo.Context) error {
	invID := c.Param("invitationID")
	_ = invID
	// TODO: mark invitation accepted, create match + chatroom, notify both users
	return c.JSON(http.StatusOK, echo.Map{
		"match_id": "stub-match-id",
		"message":  "invitation accepted",
	})
}

// POST /api/v1/invitations/:invitationID/reject  (Free tier only)
func (h *InvitationHandler) RejectInvitation(c echo.Context) error {
	invID := c.Param("invitationID")
	_ = invID
	// TODO: mark invitation dismissed; Privilege user is not notified
	return c.JSON(http.StatusOK, echo.Map{"message": "invitation rejected"})
}
