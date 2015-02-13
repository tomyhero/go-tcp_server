package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tomyhero/ore_server/client"
	"github.com/tomyhero/ore_server/context"
	"github.com/tomyhero/ore_server/server"
	"github.com/tomyhero/ore_server/util"
	"github.com/ugorji/go/codec"
	"reflect"
	"sync"
	"testing"
	"time"
)

var wg *sync.WaitGroup = &sync.WaitGroup{}

var try = 3
var countMessage = 0
var countBroadcast = 0

func TestChat(t *testing.T) {

	port, err := util.EmptyPort()
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}

	var ch = new(codec.MsgpackHandle)
	ch.MapType = reflect.TypeOf(map[string]interface{}{})
	ch.RawToString = true
	config := &server.ServerConfig{Port: port}
	sv := server.NewServer(config)
	defer sv.Shutdown()

	handlers := make([]context.IHandler, 1)
	handlers[0] = NewChatHandler()
	sv.Setup(handlers)
	go sv.Run()

	time.Sleep(100 * time.Millisecond)

	var h = new(codec.MsgpackHandle)
	h.MapType = reflect.TypeOf(map[string]interface{}{})
	h.RawToString = true

	for i := 0; i < try; i++ {
		cl := client.Client{
			CDataManager: &context.CDataManager{CodecHandle: h},
		}

		err = cl.Connect(fmt.Sprintf(":%d", port))
		if err != nil {
			fmt.Println(err)
			t.Fail()
			return
		}

		wg.Add(1)
		go ReceiveHandler(i, &cl, t)

		err = cl.Send(&context.CData{
			Header: map[string]interface{}{"CMD": "chat_login"},
			Body:   map[string]interface{}{},
		})

	}

	assert.Equal(t, 1, 1)

	wg.Wait()

	assert.Equal(t, try, countBroadcast)
	assert.NotEqual(t, 0, countMessage)

}

func ReceiveHandler(id int, cl *client.Client, t *testing.T) {
	defer func() {
		cl.Disconnect()
		wg.Done()
	}()

	for {
		res, err := cl.Receive()
		if err != nil {
			fmt.Println("ERROR", err)
			return
		}

		fmt.Println("CRes", res.GetCMD())

		switch {
		case "chat_login_res" == res.GetCMD():

			accessToken := res.Body["AUTH_ACCESS_TOKEN"]
			assert.NotNil(t, accessToken)
			hoge := "hgoe ghoegeh ogehoge ogehoego hgeo poafdij paodf japdsofaj pdoijf pdasfo daof ijadopf j"
			for i := 0; i < 10; i++ {
				hoge = hoge + hoge
			}
			hoge = fmt.Sprintf("I am alive! %d %s", id, hoge)
			fmt.Println(fmt.Sprintf("Send Message Len %d", len(hoge)))
			err := cl.Send(&context.CData{
				Header: map[string]interface{}{"CMD": "chat_Broadcast", "AUTH_ACCESS_TOKEN": accessToken},
				Body:   map[string]interface{}{"name": fmt.Sprintf("tomyhero_%d", id), "message": hoge},
			})
			if err != nil {
				fmt.Println(err)
				t.Fail()
				return
			}
		case "chat_Broadcast_res" == res.GetCMD():
			countBroadcast = countBroadcast + 1
			return
		case "chat_message" == res.GetCMD():
			fmt.Println("Receive Message Len", len(res.Body["message"].(string)))
			countMessage = countMessage + 1
		}
	}
}
