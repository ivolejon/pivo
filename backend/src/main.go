package main

import (
	"fmt"

	"github.com/ivolejon/pivo/settings"
	"github.com/ivolejon/pivo/web"
)

func main() {
	env := settings.Environment()
	server := NewHTTPServer()
	web.SetupProjectRoutes(server)
	web.SetupWebsocket(server)
	if err := server.Run(env.WebServerPort); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
