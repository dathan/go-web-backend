package routes

import (
	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/auth"
	"github.com/System-Glitch/goyave/v3/cors"
	"github.com/dathan/go-web-backend/pkg/entities"
	localauth "github.com/dathan/go-web-backend/pkg/http/services/auth"
	"github.com/dathan/go-web-backend/pkg/http/services/contacts"
	"github.com/dathan/go-web-backend/pkg/http/services/hello"
	"github.com/dathan/go-web-backend/pkg/http/services/register"
	"github.com/dathan/go-web-backend/pkg/http/services/upload"
)

// Register is very intresting. router package methods generate a new route on each call remembering the last route in something called the parent so none of the objects go out of scope
func Register(router *goyave.Router) {

	// Applying default CORS settings (allow all methods and all origins)
	// Learn more about CORS options here: https://system-glitch.github.io/goyave/guide/advanced/cors.html
	router.CORS(cors.Default())

	loggedInService := &localauth.JWTAuthenticator{}
	authenticator := auth.Middleware(&entities.User{}, loggedInService)

	// Route to register
	router.Post("/register", register.Register).Validate(register.Request)

	// Route to jwt login
	jwtRouter := router.Subrouter("/auth")
	jwtRouter.Route("POST", "/login", loggedInService.Login).Validate(localauth.LoginRequest)
	jwtRouter.Route("POST", "/refresh", loggedInService.Refresh).Validate(localauth.RefreshRequest)
	jwtRouter.Route("GET", "/google", loggedInService.GoogleLogin)
	jwtRouter.Route("GET", "/google/callback", loggedInService.GoogleAuthCallBack)

	// Route login required
	router.Get("/hello", hello.SayHi).Middleware(authenticator)

	// Route with validation
	router.Post("/echo", hello.Echo).Validate(hello.EchoRequest)

	// You must be logged in to upload a file
	router.Post("/csv/upload", upload.CSVUpload).Middleware(authenticator)
	router.Get("/contacts", contacts.List).Middleware(authenticator)
	router.Post("/contacts", contacts.Add).Middleware(authenticator).Validate(contacts.AddContactRequest)

}
