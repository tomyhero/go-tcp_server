package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tomyhero/ore_server/client"
	"github.com/tomyhero/ore_server/context"
	"github.com/tomyhero/ore_server/example/handler"
	"github.com/tomyhero/ore_server/util"
	"github.com/ugorji/go/codec"
	"reflect"
	"testing"
	"time"
)

func TestUDPServer(t *testing.T) {

	port, err := util.EmptyPort()
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}

	config := &UDPServerConfig{Port: port}

	network := &UDPNetwork{Config: config}
	sv := NewServer(network)

	assert.NotNil(t, config)

	handlers := make([]context.IHandler, 1)
	handlers[0] = handler.NewEchoHandler()
	sv.Setup(handlers)
	go sv.Run()

	time.Sleep(100 * time.Millisecond)

	var h = new(codec.MsgpackHandle)
	h.MapType = reflect.TypeOf(map[string]interface{}{})
	h.RawToString = true

	cl := client.Client{
		CDataManager: &context.CDataManager{CodecHandle: h},
	}
	err = cl.Connect(fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}

	req := &context.CData{
		Header: map[string]interface{}{"CMD": "echo_Echo", "AUTH_PLAIN_PASSWORD": "1111"},
		Body:   map[string]interface{}{"text": "Hello World\n"},
	}
	err = cl.Send(req)
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	res, err := cl.Receive()
	assert.Equal(t, "echo_Echo_res", res.GetCMD(), "response command ok")
}
