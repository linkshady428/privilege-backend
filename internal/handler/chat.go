package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ChatHandler struct{}

func NewChatHandler() *ChatHandler { return &ChatHandler{} }

type Message struct {
	MessageID string `json:"message_id"`
	SenderID  string `json:"sender_id"`
	Body      string `json:"body"`
	SentAt    string `json:"sent_at"`
}

type SendMessageRequest struct {
	Body string `json:"body"`
}

// GET /api/v1/chats/:matchID/messages
// Returns last 100 messages; supports cursor-based pagination via ?before=<message_id>
func (h *ChatHandler) GetMessages(c echo.Context) error {
	matchID := c.Param("matchID")
	_ = matchID
	// TODO: query messages from DB ordered by sent_at DESC, limit 100
	messages := []Message{
		{MessageID: "stub-msg-1", SenderID: "partner-1", Body: "Hey!", SentAt: "2024-01-01T00:00:00Z"},
	}
	return c.JSON(http.StatusOK, echo.Map{"messages": messages})
}

// POST /api/v1/chats/:matchID/messages
// REST fallback for sending a message (primary path is WebSocket).
func (h *ChatHandler) SendMessage(c echo.Context) error {
	matchID := c.Param("matchID")
	_ = matchID
	var req SendMessageRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// TODO: persist message, fan out via WebSocket hub
	return c.JSON(http.StatusCreated, Message{
		MessageID: "stub-msg-new",
		SenderID:  "dev-user-id",
		Body:      req.Body,
		SentAt:    "2024-01-01T00:00:01Z",
	})
}
