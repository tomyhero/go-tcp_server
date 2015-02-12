package server

import (
	"fmt"
	"github.com/tomyhero/ore_server/context"
	"github.com/ugorji/go/codec"
	"io"
	"net"
	"os"
	"reflect"
	"sync"
	"time"
)

type Server struct {
	Port             int                      // listen port number 
	dispatcher       *Dispatcher              // hold handlers and dispatch to it
	gstore           map[string]interface{}   // global data storage
	quit             chan bool                // chan for quit trigger.
	waitQuitGroup    *sync.WaitGroup          // wait for graceful stop
	connStore        map[net.Conn]interface{} // store all connection related data
	CodecHandle      codec.Handle
	DeadLineMillisec time.Duration
}

func (s *Server) Setup(handlers []context.IHandler) {
	s.quit = make(chan bool)
	s.waitQuitGroup = &sync.WaitGroup{}
	s.dispatcher = NewDispatcher(handlers)
	s.gstore = map[string]interface{}{}
	s.connStore = map[net.Conn]interface{}{}
	s.waitQuitGroup.Add(1)

	if s.DeadLineMillisec == 0 {
		s.DeadLineMillisec = 1000 * 1 // 1 Sec
	}

}

func (s *Server) Run() {
	defer s.waitQuitGroup.Done()
	fmt.Println("Starting Server")
	defer fmt.Println("Stoping Server")

	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ln, err := net.ListenTCP("tcp", laddr)

	defer func() {
		ln.Close()
		s.dispatcher.HookDestroy(s.gstore)
	}()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	s.dispatcher.HookInitialize(s.gstore)

	for {
		select {
		case <-s.quit:
			//fmt.Println("stopping listening on", ln.Addr())
			return
		default:
		}

		ln.SetDeadline(time.Now().Add(s.DeadLineMillisec * time.Millisecond))
		conn, err := ln.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			fmt.Println("Error Accepting", err.Error())
			os.Exit(1)
		}
		s.waitQuitGroup.Add(1)
		s.connStore[conn] = map[string]interface{}{}

		cm := &context.CDataManager{CodecHandle: s.CodecHandle}
		go s.handle(s.dispatcher, cm, conn)
	}
}

func (s *Server) Shutdown() {
	close(s.quit)
	//fmt.Println("Waiting Shutdown...")
	s.waitQuitGroup.Wait()
}

// Handles incoming requests.
func (s *Server) handle(dispatcher *Dispatcher, cm *context.CDataManager, conn net.Conn) {
	defer func() {
		conn.Close()
		delete(s.connStore, conn)
	}()
	defer s.waitQuitGroup.Done()

	for {
		select {
		case <-s.quit:
			//fmt.Println("disconnecting", conn.RemoteAddr())
			return
		default:
		}

		//	fmt.Println("start")
		conn.SetDeadline(time.Now().Add(s.DeadLineMillisec * time.Millisecond))
		data, err := cm.Receive(conn)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			if err == io.EOF {
				fmt.Println("client dissconected")
				break
			} else {
				fmt.Println("receive cdata", err)
				break
			}
		}

		c, err := context.NewContext(conn, cm, s.gstore, data, s.connStore)
		if err != nil {
			fmt.Println("create context", err)
			break
		}

		loginAction, onLogin := dispatcher.LoginActions[c.Req.GetCMD()]
		// do login
		if onLogin {
			c.SetupMyStore()
			//fmt.Println(c.Req.GetCMD(), c.Req.Header, dispatcher.LoginActions, loginAction)

			ok := loginAction.Call([]reflect.Value{reflect.ValueOf(c)})[0].Bool()
			if ok {
				c.Res.Header["STATUS"] = context.STATUS_OK
			} else {
				c.Res.Header["STATUS"] = context.STATUS_NOT_OK
			}
			// do auth action 
		} else {

			action, find := dispatcher.Actions[c.Req.GetCMD()]
			//fmt.Println(c.Req.GetCMD(), c.Req.Header, dispatcher.Actions, action, find)
			if find {
				c.SetupMyStore()

				if dispatcher.ExecAuth(c, c.Req.GetCMD()) {

					// BEFORE_EXECUTE
					dispatcher.BeforeExecute(c, c.Req.GetCMD())

					action.Call([]reflect.Value{reflect.ValueOf(c)})
					c.Res.Header["STATUS"] = context.STATUS_OK

					// AFTER_EXECUTE
					dispatcher.AfterExecute(c, c.Req.GetCMD())

				} else {
					c.Res.Header["STATUS"] = context.STATUS_FORBIDDEN
				}
			} else {
				c.Res.Header["STATUS"] = context.STATUS_COMMAND_NOT_FOUND
			}
		}

		err = cm.Send(conn, c.Res.GetData())
		if err != nil {
			fmt.Println("send fail", err)
			break
		}
		fmt.Println("end")
	}
}
