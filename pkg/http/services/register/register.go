package register

import (
	"errors"
	"net/http"
	"strings"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/auth"
	"github.com/System-Glitch/goyave/v3/database"
	"github.com/System-Glitch/goyave/v3/lang"
	userentity "github.com/dathan/go-web-backend/pkg/entities/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// This is a non auth endpoint, Given the request add to the database and return the nextUrl (login)
func Register(response *goyave.Response, request *goyave.Request) {

	username := request.String("username")
	email := request.String("email")
	password := request.String("password")

	user := userentity.User{}

	columns := auth.FindColumns(user, "username", "email", "password")

	result := database.GetConnection().Where(columns[0].Name+" = ?", username).First(user)
	notFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if result.Error != nil && !notFound {
		response.JSON(http.StatusUnprocessableEntity, map[string]string{"validationError": lang.Get(request.Lang, "register.account_exists")})
		return
	}

	result = database.GetConnection().Where(columns[1].Name+" = ?", email).First(user)
	notFound = errors.Is(result.Error, gorm.ErrRecordNotFound)

	if result.Error != nil && !notFound {
		response.JSON(http.StatusUnprocessableEntity, map[string]string{"validationError": lang.Get(request.Lang, "register.account_exists")})
		return
	}

	// todo: abstract the framework out
	// todo: abstract the datbase layer from the framework
	user.Email = email
	hpass, err := HashPassword(password)

	if err != nil {
		response.JSON(http.StatusUnprocessableEntity, map[string]string{"validationError": lang.Get(request.Lang, "register.password")})
		return
	}

	user.Password = hpass
	user.UserName = strings.ToLower(username)
	user.Email = email
	user.Birthday = nil // there is an upstream bug with 0000-00-00 date defults, set this to nil

	result = database.GetConnection().Create(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	response.JSON(http.StatusOK, user)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
