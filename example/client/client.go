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

	//for {
	cdata, err := cl.Receive()
	fmt.Println(cdata, err)
	//}

}
