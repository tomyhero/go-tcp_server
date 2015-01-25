package main

import (
	"github.com/tomyhero/ore_server/context"
	"github.com/tomyhero/ore_server/sample"
	"github.com/tomyhero/ore_server/server"
)

func main() {
	sv := server.Server{Port: 8080}
	handlers := make([]context.IHandler, 1)
	handlers[0] = sample.NewEchoHandler()
	sv.Setup(handlers)
	sv.Run()
}
