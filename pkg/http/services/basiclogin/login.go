package basiclogin

import (
	"github.com/System-Glitch/goyave/v3"
	localauth "github.com/dathan/go-web-backend/pkg/auth"
)

// Login takes in a response and a request for the framework used. TODO:move this into generic logic
// ## generic algorithim:
//  * verify the format of the username and password
//  * check the database for that combo
//  * if the entry does not exist throw error TODO:add a rate-limiter check at some level
//	* else send back the authorization information
// ### Note maybe add some sort of cookie store? or create LoginWithCookie and create a verify?
func Login(response *goyave.Response, request *goyave.Request) {

	auth := localauth.JWTAuthenticator{}
	auth.Login(response, request)
}
