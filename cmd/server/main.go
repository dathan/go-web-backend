package main

import (
	"os"

	"github.com/dathan/go-web-backend/pkg/server"
)

func main() {

	if err := server.New().Start(); err != nil {
		os.Exit(-1)
	}
}
