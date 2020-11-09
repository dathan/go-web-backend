package register

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/auth"
	"github.com/System-Glitch/goyave/v3/database"
	"github.com/System-Glitch/goyave/v3/lang"
	"github.com/dathan/go-web-backend/pkg/entities"
	localresponse "github.com/dathan/go-web-backend/pkg/http/response"
	userservice "github.com/dathan/go-web-backend/pkg/http/services/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// This is a non auth endpoint, Given the request add to the database and return the nextUrl (login)
func Register(response *goyave.Response, request *goyave.Request) {

	username := request.String("username")
	email := request.String("email")
	password := request.String("password")

	user := &entities.User{}
	var err error
	if userservice.Exists("email", email, user) {
		validationError("register.account_exists_username", request, response)
		return
	}

	if userservice.Exists("username", username, user) {
		validationError("register.account_exists_username", request, response)
		return
	}

	if user, err = RegisterUser(username, email, password); err != nil {
		validationError(err.Error(), request, response)
		return
	}

	resp := localresponse.NewResponse(true)
	resp.User = user
	response.JSON(http.StatusOK, resp)
}

func validationError(langKey string, request *goyave.Request, response *goyave.Response) {
	resp := localresponse.NewResponse(false)
	resp.ErrorMessage = fmt.Sprintf("validationError: %s", lang.Get(request.Lang, langKey))
	response.JSON(http.StatusUnprocessableEntity, resp)
	return
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RegisterUser(username, email, password string) (*entities.User, error) {
	user := entities.User{}

	columns := auth.FindColumns(user, "username", "email", "password")

	// spew.Dump(columns)
	// TODO: Refactor this block into a account_exists_service

	result := database.GetConnection().Where(columns[0].Name+" = ?", username).First(&user)
	notFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if result.Error != nil && !notFound || user.ID > 0 {
		return &user, result.Error

	}

	result = database.GetConnection().Where(columns[1].Name+" = ?", email).First(&user)
	notFound = errors.Is(result.Error, gorm.ErrRecordNotFound)

	if result.Error != nil && !notFound || user.ID > 0 {
		return &user, result.Error
	}

	// Hash the Password and store it.
	hpass, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	// set the user
	user.Password = hpass
	user.UserName = strings.ToLower(username)
	user.Email = email
	user.Birthday = nil // there is an upstream bug with 0000-00-00 date defults, set this to nil

	result = database.GetConnection().Create(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil

}
