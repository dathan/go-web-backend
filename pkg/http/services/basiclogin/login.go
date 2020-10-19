package basiclogin

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/auth"
	"github.com/System-Glitch/goyave/v3/database"
	"github.com/System-Glitch/goyave/v3/lang"
	"github.com/dathan/go-web-backend/pkg/entities"
	userentity "github.com/dathan/go-web-backend/pkg/entities/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Login takes in a response and a request for the framework used. TODO:move this into generic logic
// ## generic algorithim:
//  * verify the format of the username and password
//  * check the database for that combo
//  * if the entry does not exist throw error TODO:add a rate-limiter check at some level
//	* else send back the authorization information
// ### Note maybe add some sort of cookie store? or create LoginWithCookie and create a verify?
func Login(response *goyave.Response, request *goyave.Request) {

	user := userentity.User{}
	username := request.String("username")
	columns := auth.FindColumns(user, "username", "password")
	resp := entities.NewResponse(false)

	result := database.GetConnection().Where(columns[0].Name+" = ?", username).First(&user)
	notFound := errors.Is(result.Error, gorm.ErrRecordNotFound)

	if result.Error != nil && !notFound {
		resp.ErrorMessage = result.Error.Error()
		response.JSON(http.StatusNotFound, resp)
		return
	}

	pass := reflect.Indirect(reflect.ValueOf(user)).FieldByName(columns[1].Field.Name)
	if notFound || bcrypt.CompareHashAndPassword([]byte(pass.String()), []byte(request.String("password"))) != nil {
		resp.ErrorMessage = fmt.Sprintf("validationError: %s", lang.Get(request.Lang, "auth.invalid-credentials"))
		response.JSON(http.StatusUnprocessableEntity, resp)
		return
	}

	token, err := auth.GenerateToken(username)
	if err != nil {
		resp.ErrorMessage = err.Error()
		response.JSON(http.StatusNotFound, resp)
		return
	}

	resp.OK = true
	resp.Token = token
	resp.User = &user
	response.JSON(http.StatusOK, resp)
	return
}
