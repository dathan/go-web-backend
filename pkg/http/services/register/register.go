package register

import (
	"fmt"
	"net/http"

	"github.com/System-Glitch/goyave/v3"
)

func Register(response *goyave.Response, request *goyave.Request) {

	response.String(http.StatusOK, request.String("text"))

	username := request.String("username")
	email := request.String("email")
	password := request.String("password")

	// todo: abstract the framework out
	// todo: abstract the datbase layer from the framework

	response.String(http.StatusMethodNotAllowed, fmt.Sprintf("NOT IMPLEMENTED: Username: %s Email: %s Password %s", username, email, password))

}
