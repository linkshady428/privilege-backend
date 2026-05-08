package handler

import (
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	mw "github.com/privilege/backend/internal/middleware"
)

type UserHandler struct {
	db *pgxpool.Pool
}

func NewUserHandler(pool *pgxpool.Pool) *UserHandler { return &UserHandler{db: pool} }

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
	ctx := c.Request().Context()

	var p UserProfile
	var birthdate time.Time
	err := h.db.QueryRow(ctx, `
		SELECT id, name, birthdate,
		       COALESCE(bio, ''), COALESCE(job, ''),
		       sex::text, COALESCE(relationship_status::text, ''),
		       COALESCE(height_cm, 0)::int, COALESCE(weight_kg, 0)::int,
		       lifestyle_tags, tier::text
		FROM users WHERE id = $1 AND deleted_at IS NULL`, userID,
	).Scan(
		&p.ID, &p.Name, &birthdate,
		&p.Bio, &p.Job,
		&p.Sex, &p.RelationshipStatus,
		&p.HeightCm, &p.WeightKg,
		&p.LifestyleTags, &p.Tier,
	)
	if err == pgx.ErrNoRows {
		return echo.NewHTTPError(http.StatusNotFound, "user not found")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "server error")
	}

	p.Age = int(time.Since(birthdate).Hours() / (24 * 365.25))
	if p.LifestyleTags == nil {
		p.LifestyleTags = []string{}
	}

	rows, err := h.db.Query(ctx,
		`SELECT url FROM photos WHERE user_id = $1 ORDER BY position`, userID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var url string
			rows.Scan(&url)
			p.Photos = append(p.Photos, url)
		}
	}
	if p.Photos == nil {
		p.Photos = []string{}
	}

	return c.JSON(http.StatusOK, p)
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
	_ = c.Param("photoID")
	// TODO: remove from Cloudflare Images and DB
	return c.JSON(http.StatusOK, echo.Map{"message": "photo deleted"})
}

// PUT /api/v1/users/me/location
func (h *UserHandler) UpdateLocation(c echo.Context) error {
	var req UpdateLocationRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// TODO: reverse-geocode to city, store latitude/longitude + geohash
	return c.JSON(http.StatusOK, echo.Map{"message": "location updated"})
}
