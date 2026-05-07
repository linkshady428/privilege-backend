package router

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/privilege/backend/internal/config"
	"github.com/privilege/backend/internal/handler"
	mw "github.com/privilege/backend/internal/middleware"
)

func Register(e *echo.Echo, cfg *config.Config, pool *pgxpool.Pool) {
	e.GET("/health", handler.Health)

	v1 := e.Group("/api/v1")

	// --- Public: Auth ---
	auth := handler.NewAuthHandler(pool, cfg.JWTSecret)
	v1.POST("/auth/register", auth.Register)
	v1.POST("/auth/login", auth.Login)
	v1.POST("/auth/refresh", auth.Refresh)
	v1.POST("/auth/logout", auth.Logout)

	// --- Protected ---
	protected := v1.Group("", mw.JWT(cfg))

	// Users
	users := handler.NewUserHandler(pool)
	protected.GET("/users/me", users.GetMe)
	protected.PUT("/users/me", users.UpdateMe)
	protected.DELETE("/users/me", users.DeleteMe)
	protected.POST("/users/me/photos", users.UploadPhoto)
	protected.DELETE("/users/me/photos/:photoID", users.DeletePhoto)
	protected.PUT("/users/me/location", users.UpdateLocation)

	// Feed
	feed := handler.NewFeedHandler(pool)
	protected.GET("/feed", feed.GetFeed)
	protected.POST("/feed/pass/:userID", feed.PassUser, mw.RequireTier("privilege"))

	// Invitations
	invitations := handler.NewInvitationHandler()
	protected.POST("/invitations", invitations.SendInvitation, mw.RequireTier("privilege"))
	protected.POST("/invitations/:invitationID/accept", invitations.AcceptInvitation, mw.RequireTier("free"))
	protected.POST("/invitations/:invitationID/reject", invitations.RejectInvitation, mw.RequireTier("free"))

	// Matches
	matches := handler.NewMatchHandler()
	protected.GET("/matches", matches.ListMatches)

	// Chat (REST + WebSocket)
	chat := handler.NewChatHandler()
	protected.GET("/chats/:matchID/messages", chat.GetMessages)
	protected.POST("/chats/:matchID/messages", chat.SendMessage)

	// WebSocket — path intentionally outside /api/v1
	ws := handler.NewWSHandler()
	e.GET("/ws/chat/:matchID", ws.Connect, mw.JWT(cfg))

	// Safety
	safety := handler.NewSafetyHandler()
	protected.POST("/users/:userID/block", safety.BlockUser)
	protected.POST("/users/:userID/report", safety.ReportUser)
}
