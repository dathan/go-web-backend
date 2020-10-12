package server

import (
	"github.com/System-Glitch/goyave/v3"
	_ "github.com/System-Glitch/goyave/v3/database/dialect/mysql" // this loads the driver and call the package init
	_ "github.com/dathan/go-web-backend/pkg/entities/user"        // this registers the model at startup to create the instances.
	"github.com/dathan/go-web-backend/pkg/http/routes"
)

//interface for Starting the server
type Server interface {
	Start() error
}

//server struct stuv that implements the interface
type server struct {
}

func New() Server {

	return &server{}
}

// start the server
func (s *server) Start() error {

	// This is the entry point of your application.
	var err error
	if err = goyave.Start(routes.Register); err != nil {
		return err
	}

	return nil

}
