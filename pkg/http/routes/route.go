package routes

import (
	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/cors"
	"github.com/dathan/go-web-backend/pkg/http/services/basiclogin"
	"github.com/dathan/go-web-backend/pkg/http/services/hello"
	"github.com/dathan/go-web-backend/pkg/http/services/register"
)

// Register is very intresting. router package methods generate a new route on each call remembering the last route in something called the parent so none of the objects go out of scope
func Register(router *goyave.Router) {

	// Applying default CORS settings (allow all methods and all origins)
	// Learn more about CORS options here: https://system-glitch.github.io/goyave/guide/advanced/cors.html
	router.CORS(cors.Default())

	// Register your routes here

	// Route without validation
	router.Get("/hello", hello.SayHi)

	// Route with validation
	router.Post("/echo", hello.Echo).Validate(hello.EchoRequest)

	// Route to register
	router.Post("/register", register.Register).Validate(register.Request)

	// Route to login
	router.Post("/auth/login", basiclogin.Login).Validate(basiclogin.Request)

	//router.Post("/auth/google/callback", idp.Google).Validate(idp.Google)

	//auth.Middleware(&model.User{}, &auth.JWTAuthenticator{})
	//router.Middleware()
}
