package main

import (
	"github.com/tomyhero/go-tcp_server/context"
	"github.com/tomyhero/go-tcp_server/example/handler"
	"github.com/tomyhero/go-tcp_server/server"
)

func main() {
	config := &server.ServerConfig{Port: 8080}

	sv := server.NewServer(config)
	handlers := make([]context.IHandler, 1)
	handlers[0] = handler.NewEchoHandler()
	sv.Setup(handlers)
	sv.Run()
}
