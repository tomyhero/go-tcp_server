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

	// when out of for loop, close the connection.
	defer conn.Close()

	for {
		cdata := CData{SerializorType: SERIALIZOR_TYPE_MESSAGE_PACK}
		data, err := cdata.Receive(conn)
		if err != nil {
			fmt.Println("receive cdata", err)
			break
		}

		c, err := context.NewContext(data)
		if err != nil {
			fmt.Println("create context", err)
			break
		}

		// Authorize or Login
		// HOOK_BEFORE

		action, find := dispatcher.Actions[c.Req.GetCMD()]
		fmt.Println(c.Req.GetCMD(), c.Req.Header, dispatcher.Actions, action, find)
		if find {
			action.Call([]reflect.Value{reflect.ValueOf(c)})
		} else {

		}

		// HOOK_AFTER

	}
}
