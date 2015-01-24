package main

import (
	"github.com/tomyhero/ore_server/server"
)

func main() {
	sv := server.Server{Port: 8080}
	sv.Run()
}
