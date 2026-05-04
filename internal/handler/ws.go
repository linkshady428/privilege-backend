package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WSHandler struct{}

func NewWSHandler() *WSHandler { return &WSHandler{} }

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// GET /ws/chat/:matchID
// Upgrades to WebSocket for real-time messaging.
// Currently echoes messages back; TODO: integrate with hub + DB persistence.
func (h *WSHandler) Connect(c echo.Context) error {
	matchID := c.Param("matchID")
	_ = matchID

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		if err := conn.WriteMessage(mt, msg); err != nil {
			break
		}
	}
	return nil
}
