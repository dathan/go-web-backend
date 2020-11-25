package entities

import (
	"github.com/System-Glitch/goyave/v3/database"
	"github.com/jinzhu/gorm"
)

// A model is a structure reflecting a database table structure. An instance of a model
// is a single database record. Each model is defined in its own file inside the database/models directory.
// Models are usually just normal Golang structs, basic Go types, or pointers of them.
// "sql.Scanner" and "driver.Valuer" interfaces are also supported.

func init() {
	// All models should be registered in an "init()" function inside their model file.
	database.RegisterModel(&Contacts_Paths{})
}

//contacts_path holds a pointer to a file which exists on s3 or local disk to save space
type Contacts_Paths struct {
	gorm.Model
	OwnerID      uint   `gorm:"index" json:"owner_id"`
	FileLocation string `gorm:"type:varchar(1024)"`
	CSVData      string `gorm:"type:blob"`
}
