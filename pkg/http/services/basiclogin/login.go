package basiclogin

import (
	"fmt"
	"net/http"

	"github.com/System-Glitch/goyave/v3"
)

// Login takes in a response and a request for the framework used. TODO:move this into generic logic
// ## generic algorithim:
//  * verify the format of the username and password
//  * check the database for that combo
//  * if the entry does not exist throw error TODO:add a rate-limiter check at some level
//	* else send back the authorization information
// ### Note maybe add some sort of cookie store? or create LoginWithCookie and create a verify?
func Login(response *goyave.Response, request *goyave.Request) {
	// Note: Depending on the request you may get back a generic response
	//{
	//	"error": "Field \"email\" is not a string"
	//}
	// This is a default error thrown as a result of the package, might want to intercept it and process it.
	response.String(http.StatusNotImplemented, fmt.Sprintf("NOT IMPLEMENT: %s, %s", request.String("username"), request.String("password")))
}
