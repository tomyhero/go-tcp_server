package server

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/tomyhero/go-tcp_server/context"
	"github.com/ugorji/go/codec"
	"io"
	"net"
	"reflect"
	"sync"
	"time"
)

type ServerConfig struct {
	Port             int
	CodecHandle      codec.Handle
	DeadLineMillisec time.Duration
}

func (config *ServerConfig) LoadDefault() {
	if config.Port == 0 {
		config.Port = 8081
	}

	if config.CodecHandle == nil {
		var h = new(codec.MsgpackHandle)
		h.MapType = reflect.TypeOf(map[string]interface{}{})
		h.RawToString = true
		config.CodecHandle = h
	}

	if config.DeadLineMillisec == 0 {
		config.DeadLineMillisec = 1000 // 1 sec
	}
}

// ---

type Server struct {
	dispatcher    *Dispatcher              // hold handlers and dispatch to it
	gstore        map[string]interface{}   // global data storage
	quit          chan bool                // chan for quit trigger.
	waitQuitGroup *sync.WaitGroup          // wait for graceful stop
	connStore     map[net.Conn]interface{} // store all connection related data
	Config        *ServerConfig
}

func NewServer(config *ServerConfig) *Server {
	config.LoadDefault()
	return &Server{Config: config}
}

func (s *Server) Setup(handlers []context.IHandler) {
	s.quit = make(chan bool)
	s.waitQuitGroup = &sync.WaitGroup{}
	s.dispatcher = NewDispatcher(handlers)
	s.gstore = map[string]interface{}{}
	s.connStore = map[net.Conn]interface{}{}
	s.waitQuitGroup.Add(1)
}

func (s *Server) Run() {
	fmt.Println("Run Called")
	defer func() {
		glog.Info("End of Run")
		s.waitQuitGroup.Done()
	}()

	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", s.Config.Port))

	if err != nil {
		glog.Fatalf("ResolveTCPAddr Fail: %s", err)
	}

	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		glog.Fatalf("ListenTCP Fail: %s", err)
	}

	defer func() {
		ln.Close()
		glog.Info("Listener Closed")
		s.dispatcher.HookDestroy(s.gstore)
	}()

	s.dispatcher.HookInitialize(s.gstore)

	for {
		select {
		case <-s.quit:
			glog.Info("quit Command Received")
			return
		default:
		}

		ln.SetDeadline(time.Now().Add(s.Config.DeadLineMillisec * time.Millisecond))
		conn, err := ln.Accept()

		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			glog.Fatalf("Accept Fail: %s", err)
		}

		fmt.Println("Connected")

		s.waitQuitGroup.Add(1)
		s.connStore[conn] = map[string]interface{}{}

		cm := &context.CDataManager{CodecHandle: s.Config.CodecHandle}
		go s.handle(s.dispatcher, cm, conn)
	}
}

// Handles incoming requests.
func (s *Server) handle(dispatcher *Dispatcher, cm *context.CDataManager, conn net.Conn) {
	defer func() {
		conn.Close()
		delete(s.connStore, conn)
		glog.Info("Close Connection")
	}()
	defer s.waitQuitGroup.Done()

	for {
		select {
		case <-s.quit:
			glog.Infof("quit Command Received at Handle %s", conn.RemoteAddr())
			return
		default:
		}

		conn.SetDeadline(time.Now().Add(s.Config.DeadLineMillisec * time.Millisecond))
		data, err := cm.Receive(conn)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			if err == io.EOF {
				glog.Infof("Client Dissconected :%s", conn.RemoteAddr())
				break
			} else {
				glog.Infof("Receive Data Failed %s", err)
				break
			}
		}

		_, has_header := data["H"]
		if has_header == false {
			continue
		}

		fmt.Println("data", data)

		c, err := context.NewContext(conn, cm, s.gstore, data, s.connStore)
		if err != nil {
			glog.Warningf("Creating Context Failed %s", err)
			break
		}

		loginAction, onLogin := dispatcher.LoginActions[c.Req.GetCMD()]
		// do login
		if onLogin {
			c.SetupMyStore()
			ok := loginAction.Call([]reflect.Value{reflect.ValueOf(c)})[0].Bool()
			if ok {
				c.Res.Header["STATUS"] = context.STATUS_OK
			} else {
				c.Res.Header["STATUS"] = context.STATUS_NOT_OK
			}
			// do auth action 
		} else {
			action, find := dispatcher.Actions[c.Req.GetCMD()]
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

		if c.OnSendResponse {
			fmt.Println("sending", c.Res.GetData())
			err = cm.Send(conn, c.Res.GetData())
			fmt.Println("sent", c.Res.GetData())
			if err != nil {
				glog.Infof("Sending Data Fail %s", err)
				break
			}
		}
	}
}

func (s *Server) Shutdown() {
	close(s.quit)
	glog.Info("Waiting Shutdown....")
	s.waitQuitGroup.Wait()
}
