package server

/*
sv := Server{"Port":8080,"Handlers":{FooHandler,BooHandler}}
sv.Run()
*/

import (
	"fmt"
	"net"
	"os"
)

type Server struct {
	Port       int
	conns      []net.Conn
	ln         net.Listener
	dispatcher *Dispatcher
}

func (s *Server) Setup(handlers interface{}) {
	s.dispatcher = NewDispatcher(handlers)
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return err
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error Accepting", err.Error())
			os.Exit(1)
		}
		go handle(conn)
	}
}

// Handles incoming requests.
func handle(conn net.Conn) {
	for {
		// Make a buffer to hold incoming data.
		buf := make([]byte, 1024)
		// Read the incoming connection into the buffer.
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		fmt.Println(string(buf))
	}
}
