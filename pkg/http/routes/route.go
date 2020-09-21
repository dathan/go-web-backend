package routes

import (
	"github.com/System-Glitch/goyave/v3"
	"github.com/System-Glitch/goyave/v3/cors"
	"github.com/dathan/go-web-backend/pkg/http/services/hello"
)

func Register(router *goyave.Router) {

	// Applying default CORS settings (allow all methods and all origins)
	// Learn more about CORS options here: https://system-glitch.github.io/goyave/guide/advanced/cors.html

	router.CORS(cors.Default())

	// Register your routes here

	// Route without validation
	router.Get("/hello", hello.SayHi)

	// Route with validation
	router.Post("/echo", hello.Echo).Validate(hello.EchoRequest)
}
