module github.com/dathan/go-web-backend

go 1.15

replace (
	github.com/dathan/go-web-backend => /Users/dathan.pattishall/workspace/go-web-backend
	github.com/dathan/go-web-backend/pkg/http/services/auth => /Users/dathan.pattishall/workspace/go-web-backend/pkg/http/services/auth
)

require (
	github.com/System-Glitch/goyave v1.0.0
	github.com/System-Glitch/goyave/v3 v3.1.0
	github.com/davecgh/go-spew v1.1.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/jinzhu/gorm v1.9.11
	github.com/markbates/goth v1.65.0
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	gorm.io/gorm v1.20.1
)
