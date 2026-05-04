package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	mw "github.com/privilege/backend/internal/middleware"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler { return &UserHandler{} }

type UserProfile struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Age                int      `json:"age"`
	Bio                string   `json:"bio"`
	Photos             []string `json:"photos"`
	Job                string   `json:"job,omitempty"`
	HeightCm           int      `json:"height_cm,omitempty"`
	WeightKg           int      `json:"weight_kg,omitempty"`
	Sex                string   `json:"sex"`
	RelationshipStatus string   `json:"relationship_status,omitempty"`
	LifestyleTags      []string `json:"lifestyle_tags,omitempty"`
	Tier               string   `json:"tier"`
}

type UpdateProfileRequest struct {
	Name               string   `json:"name"`
	Bio                string   `json:"bio"`
	Job                string   `json:"job"`
	HeightCm           int      `json:"height_cm"`
	WeightKg           int      `json:"weight_kg"`
	Sex                string   `json:"sex"`
	RelationshipStatus string   `json:"relationship_status"`
	LifestyleTags      []string `json:"lifestyle_tags"`
}

type UpdateLocationRequest struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// GET /api/v1/users/me
func (h *UserHandler) GetMe(c echo.Context) error {
	userID := c.Get(mw.ContextKeyUserID).(string)
	tier := c.Get(mw.ContextKeyTier).(string)
	// TODO: fetch from DB
	return c.JSON(http.StatusOK, UserProfile{
		ID:     userID,
		Name:   "Stub User",
		Age:    25,
		Bio:    "",
		Photos: []string{},
		Tier:   tier,
	})
}

// PUT /api/v1/users/me
func (h *UserHandler) UpdateMe(c echo.Context) error {
	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// TODO: validate photo count limit by tier, persist to DB
	return c.JSON(http.StatusOK, echo.Map{"message": "profile updated"})
}

// DELETE /api/v1/users/me
func (h *UserHandler) DeleteMe(c echo.Context) error {
	// TODO: soft-delete with 30-day grace period, void all invitations, lock chats
	return c.JSON(http.StatusOK, echo.Map{"message": "account deletion scheduled"})
}

// POST /api/v1/users/me/photos
func (h *UserHandler) UploadPhoto(c echo.Context) error {
	// TODO: receive multipart file, upload to Cloudflare Images, store URL
	return c.JSON(http.StatusCreated, echo.Map{"photo_url": "https://example.com/stub-photo.jpg"})
}

// DELETE /api/v1/users/me/photos/:photoID
func (h *UserHandler) DeletePhoto(c echo.Context) error {
	photoID := c.Param("photoID")
	_ = photoID
	// TODO: remove from Cloudflare Images and DB
	return c.JSON(http.StatusOK, echo.Map{"message": "photo deleted"})
}

// PUT /api/v1/users/me/location
func (h *UserHandler) UpdateLocation(c echo.Context) error {
	var req UpdateLocationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// TODO: reverse-geocode to city, store geohash
	return c.JSON(http.StatusOK, echo.Map{"message": "location updated"})
}
