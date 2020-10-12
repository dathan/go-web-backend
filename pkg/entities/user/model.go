package user

import (
	"time"

	"github.com/System-Glitch/goyave/v3/database"
	"github.com/jinzhu/gorm"
)

// A model is a structure reflecting a database table structure. An instance of a model
// is a single database record. Each model is defined in its own file inside the database/models directory.
// Models are usually just normal Golang structs, basic Go types, or pointers of them.
// "sql.Scanner" and "driver.Valuer" interfaces are also supported.

// Learn more here: https://system-glitch.github.io/goyave/guide/basics/database.html#models

func init() {
	// All models should be registered in an "init()" function inside their model file.
	database.RegisterModel(&User{})
}

// User represents a user. Note: There is some framework blead in the annotation of the model.
// We can see this is gorm.Model but auth: seems to be a hint for the goyave framework
type User struct {
	gorm.Model
	UserName string    `gorm:"type:varchar(50);unique_index" auth:"username"`
	Email    string    `gorm:"type:varchar(100);unique_index"`
	Password string    `gorm:"type:varchar(10);unique_index" auth:"password"`
	Birthday time.Time `gorm:"index"`
}
