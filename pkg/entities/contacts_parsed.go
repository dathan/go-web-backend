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
	CreatedAt        time.Time  `json:"created_at,omitempty"`
	UpdatedAt        time.Time  `json:"-"`
	DeletedAt        *time.Time `sql:"index" json:"-"`
	OwnerID          uint       `gorm:"primaryKey;autoincrement:false" json:"owner_id"`
	Email            string     `gorm:"type:varchar(255);primaryKey" json:"email"`
	DisplayName      string     `gorm:"type:varchar(255)" json:"display_name"`
	FirstName        string     `gorm:"type:varchar(255)" json:"first_name"`
	LastName         string     `gorm:"type:varchar(255)" json:"last_name"`
	StreetAddress    string     `gorm:"type:varchar(1024)" json:"steet_address"`
	CityCode         string     `gorm:"type:varchar(3)" json:"city"`
	ZipCode          string     `gorm:"type:varchar(20)" json:"zip_code"`
	StateCode        string     `gorm:"type:varchar(3)" json:"state"`
	Birthday         *time.Time `gorm:"index" json:"birthday,omitempty"`
	FacebookID       uint64     `json:"facebook_id"`
	FacebookUsername string     `gorm:"type:varchar(255)" json:"facebook_username"`
	PhoneNumber      string     `gorm:"type:varchar(255)" json:"phone"`
}
