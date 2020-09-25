package middleware

import "github.com/System-Glitch/goyave"

// Authcheck
//
// To use this middleware, assign it to a router in "http/routes/routes.go"
//
//     router.Middleware(middleware.AuthCheck) -- a new object is created and referenced for everty router package method that returns a Router.
func AuthCheck(next goyave.Handler) goyave.Handler {
	return func(response *goyave.Response, request *goyave.Request) {

		// 1st we will check if a cookie exists
		//  - if a cookie exists with a token lets verify that token
		// else get redirected to the login page
		// the call back stores the user sets the cookie
		// - according to the jwt spec 
		next(response, request) // Pass to the next handler
	}
}
