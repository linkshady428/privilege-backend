package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type SafetyHandler struct{}

func NewSafetyHandler() *SafetyHandler { return &SafetyHandler{} }

type ReportRequest struct {
	Reason string `json:"reason"` // "spam" | "inappropriate_content" | "feels_unsafe"
}

// POST /api/v1/users/:userID/block
func (h *SafetyHandler) BlockUser(c echo.Context) error {
	targetID := c.Param("userID")
	_ = targetID
	// TODO: log block action; post-MVP: enforce mutual hide + lock chat
	return c.JSON(http.StatusOK, echo.Map{"message": "user blocked"})
}

// POST /api/v1/users/:userID/report
func (h *SafetyHandler) ReportUser(c echo.Context) error {
	targetID := c.Param("userID")
	_ = targetID
	var req ReportRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// TODO: insert into reports table; post-MVP: admin review queue
	return c.JSON(http.StatusOK, echo.Map{"message": "report submitted"})
}
