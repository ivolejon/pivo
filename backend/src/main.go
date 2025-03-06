package main

import (
	"fmt"

	"github.com/ivolejon/pivo/config"
	"github.com/ivolejon/pivo/web"
)

func main() {
	env := config.Environment()
	server := NewHTTPServer()
	web.SetupDefaultRoutes(server)
	if err := server.Run(env.WebServerPort); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
