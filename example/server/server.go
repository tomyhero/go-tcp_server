package main

import (
	"github.com/tomyhero/ore_server/context"
	"github.com/tomyhero/ore_server/example/handler"
	"github.com/tomyhero/ore_server/server"
)

func main() {
	config := &server.ServerConfig{Port: 8080}

	sv := server.NewServer(config)
	handlers := make([]context.IHandler, 1)
	handlers[0] = handler.NewEchoHandler()
	sv.Setup(handlers)
	sv.Run()
}
