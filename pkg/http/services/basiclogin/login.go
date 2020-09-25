package basiclogin

import (
	"fmt"
	"net/http"

	"github.com/System-Glitch/goyave/v3"
)

func SiteLogin(response *goyave.Response, request *goyave.Request) {
	response.String(http.StatusNotImplemented, fmt.Sprintf("NOT IMPLEMENT: %s, %s, %s", request.String("username"), request.String("password"), request.String("email")))
}
