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

	// spew.Dump(columns)
	// TODO: Refactor this block into a account_exists_service
	result := database.GetConnection().Where(columns[0].Name+" = ?", username).First(&user)
	notFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if result.Error != nil && !notFound || user.ID > 0 {
		validationError("register.account_exists_username", request, response)
		return

	}

	result = database.GetConnection().Where(columns[1].Name+" = ?", email).First(&user)
	notFound = errors.Is(result.Error, gorm.ErrRecordNotFound)

	if result.Error != nil && !notFound || user.ID > 0 {
		validationError("register.account_exists_email", request, response)
		return
	}

	// Hash the Password and store it.
	hpass, err := HashPassword(password)
	if err != nil {
		validationError("register.password", request, response)
		return
	}

	// set the user
	user.Password = hpass
	user.UserName = strings.ToLower(username)
	user.Email = email
	user.Birthday = nil // there is an upstream bug with 0000-00-00 date defults, set this to nil

	result = database.GetConnection().Create(&user)

	if result.Error != nil {
		panic(result.Error)
	}

	resp := entities.NewResponse(true)
	resp.User = &user
	response.JSON(http.StatusOK, resp)
}

func validationError(langKey string, request *goyave.Request, response *goyave.Response) {
	resp := entities.NewResponse(false)
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
