go-tcp_server - This is TCP Stream Server Application Framework. 
========================================================

もくてき
=======

ハンドラーを書くだけで、簡単に TCP Stream Server を作れるのを目的としています

Example
=======

ハンドラー
```go
package handler

import (
	//"fmt"
	"github.com/tomyhero/go-tcp_server/authorizer"
	"github.com/tomyhero/go-tcp_server/context"
)

type EchoHandler struct {
	Authorizer context.IAuthorizer
}

func (h *EchoHandler) Prefix() string {
	return "echo"
}
func (h *EchoHandler) GetAuthorizer() context.IAuthorizer {
	return h.Authorizer
}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{Authorizer: authorizer.PlainPassword{Password: "1111"}}
}

func (h *EchoHandler) HookInitialize(g map[string]interface{}, myStore map[string]interface{}) {
	myStore["num"] = 0
}
func (h *EchoHandler) HookDestroy(g map[string]interface{}, myStore map[string]interface{}) {

}

func (h *EchoHandler) HookBeforeExecute(c *context.Context) {
	//fmt.Println("Called BeforeExecuteHandler", c.MyStore, c.GStore["echo"])
	myStore := c.MyStore()
	myStore["num"] = myStore["num"].(int) + 1
}

func (h *EchoHandler) HookAfterExecute(c *context.Context) {
	//fmt.Println("Called AfterExecuteHandler")
}

func (h *EchoHandler) ActionEcho(c *context.Context) (*context.Context, error) {
	c.Res.Body = c.Req.Body
	//fmt.Println(c, "Echo Echo!", c.MyStore()["num"].(int))
	return c, nil
}

```


```go
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

```

```go
package main

import (
	"fmt"
	"github.com/tomyhero/go-tcp_server/client"
	"github.com/tomyhero/go-tcp_server/context"
	"github.com/ugorji/go/codec"
	"reflect"
)

func main() {

	var h = new(codec.MsgpackHandle)
	h.MapType = reflect.TypeOf(map[string]interface{}{})
	h.RawToString = true
	cl := client.Client{
		CDataManager: &context.CDataManager{CodecHandle: h},
	}
	defer cl.Disconnect()

	err := cl.Connect(":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = cl.Send(&context.CData{
		Header: map[string]interface{}{"CMD": "echo_Echo", "AUTH_PLAIN_PASSWORD": "1111"},
		Body:   map[string]interface{}{"text": "Hello World\n"},
	})

	cdata, err := cl.Receive()

	fmt.Println(cdata, err)

}

```
