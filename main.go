package main

import (
	app "github.com/taufikardiyan28/chat/app"
)

func main() {
	cfg := &app.Config{
		Port: 8080,
	}

	server := app.Server{Config: cfg}
	server.Start()
}
