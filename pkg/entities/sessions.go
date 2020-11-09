package entities

import (
	"time"

	"github.com/System-Glitch/goyave/v3/database"
)

func init() {
	// All models should be registered in an "init()" function inside their model file.
	database.RegisterModel(&Session{})
}

// User represents a user. Note: There is some framework blead in the annotation of the model.
// We can see this is gorm.Model but auth: seems to be a hint for the goyave framework
type Session struct {
	ID                 uint64 `gorm:"primary_key"`
	UserID             uint   `gorm:"index:user_id_created_at"`
	Provider           string
	ProviderID         string
	AccessToken        string
	AccessTokenSecret  string
	RefreshToken       string
	RefreshTokenSecret string
	ExpiresAt          time.Time
	CreatedAt          time.Time `sql:"DEFAULT:'current_timestamp'" gorm:"index:user_id_created_at"`
}
