package user

import (
	"errors"

	"github.com/System-Glitch/goyave/v3/database"
	"github.com/dathan/go-web-backend/pkg/entities"
	"gorm.io/gorm"
)

//if user.Exists == true; side effect user filled
func Exists(field string, fieldValue string, user *entities.User) bool {
	result := database.GetConnection().Where(field+" = ?", fieldValue).First(user)
	if len(user.Email) < 5 { //a@a.a
		return false
	}

	notFound := errors.Is(result.Error, gorm.ErrRecordNotFound)
	return notFound
}
