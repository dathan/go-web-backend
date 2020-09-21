# Purpose

This will be the location of the various Auth Modules the goal is to check the login cookie for the credentials to verify in the db, then verify their token for twitter, google

```
package middleware

import "github.com/System-Glitch/goyave/v3"

// Middleware are handlers executed before the controller handler.
// They are a convenient way to filter, intercept or alter HTTP requests entering your application.
// Learn more here: https://system-glitch.github.io/goyave/guide/basics/middleware.html

// MyCustomMiddleware is an example middleware.
//
// To use this middleware, assign it to a router in "http/routes/routes.go"
//
//     router.Middleware(middleware.MyCustomMiddleware)
func MyCustomMiddleware(next goyave.Handler) goyave.Handler {
	return func(response *goyave.Response, request *goyave.Request) {
		// Do something
		next(response, request) // Pass to the next handler
	}
}
```