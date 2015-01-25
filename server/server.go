package server

/*
sv := Server{"Port":8080,"Handlers":{FooHandler,BooHandler}}
sv.Run()
*/

import (
	//"bytes"
	"fmt"
	"github.com/tomyhero/ore_server/context"
	"net"
	"os"
	"reflect"
)

type Server struct {
	Port       int
	conns      []net.Conn
	ln         net.Listener
	dispatcher *Dispatcher
}

func (s *Server) Setup(handlers []context.IHandler) {
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
		go handle(s.dispatcher, conn)
	}
}

// Handles incoming requests.
func handle(dispatcher *Dispatcher, conn net.Conn) {
	for {
		b := make([]byte, 1024)
		_, err := conn.Read(b)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		fmt.Println(string(b))

		action, find := dispatcher.Actions["echo_Test"]

		if find {
			action.Call([]reflect.Value{})
		} else {

		}

		/*
			buf := bytes.NewBuffer(b)
			c, err := NewContext(buf)
			if err != nil {
				fmt.Println("create context", err)
			}
		*/

	}
}
