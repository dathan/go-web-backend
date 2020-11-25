package entities

import (
	"time"

	"github.com/System-Glitch/goyave/v3/database"
)

func init() {
	// All models should be registered in an "init()" function inside their model file.
	database.RegisterModel(&Contacts_Parsed{})
}

//contacts_path holds a pointer to a file which exists on s3 or local disk to save space
type Contacts_Parsed struct {
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time `sql:"index"`
	OwnerID          uint       `gorm:"primaryKey;autoincrement:false" json:"owner_id"`
	Email            string     `gorm:"type:varchar(255);primaryKey"`
	DisplayName      string     `gorm:"type:varchar(255)"`
	FirstName        string     `gorm:"type:varchar(255)"`
	LastName         string     `gorm:"type:varchar(255)"`
	StreetAddress    string     `gorm:"type:varchar(1024)"`
	CityCode         string     `gorm:"type:varchar(3)"`
	ZipCode          string     `gorm:"type:varchar(20)"`
	StateCode        string     `gorm:"type:varchar(3)"`
	Birthday         *time.Time `gorm:"index" json:"birthday,omitempty"`
	FacebookID       uint64
	FacebookUsername string `gorm:"type:varchar(255)"`
	PhoneNumber      string `gorm:"type:varchar(255)"`
}
