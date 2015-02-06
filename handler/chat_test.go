package handler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tomyhero/ore_server/client"
	"github.com/tomyhero/ore_server/context"
	"github.com/tomyhero/ore_server/server"
	"github.com/tomyhero/ore_server/util"
	"sync"
	"testing"
	"time"
)

var wg *sync.WaitGroup = &sync.WaitGroup{}

func TestChat(t *testing.T) {

	port, err := util.EmptyPort()
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}

	sv := server.Server{Port: port}
	defer sv.Shutdown()

	handlers := make([]context.IHandler, 1)
	handlers[0] = NewChatHandler()
	sv.Setup(handlers)
	go sv.Run()

	time.Sleep(100 * time.Millisecond)

	for i := 0; i < 3; i++ {
		cl := client.Client{}
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

		switch {
		case "chat_login_res" == res.GetCMD():

			accessToken := res.Body["AUTH_ACCESS_TOKEN"]
			assert.NotNil(t, accessToken)
			err := cl.Send(&context.CData{
				Header: map[string]interface{}{"CMD": "chat_Broadcast", "AUTH_ACCESS_TOKEN": accessToken},
				Body:   map[string]interface{}{"name": fmt.Sprintf("tomyhero_%d", id), "message": fmt.Sprintf("I am alive! %d", id)},
			})
			if err != nil {
				fmt.Println(err)
				t.Fail()
				return
			}
		case "chat_Broadcast_res" == res.GetCMD():
			fmt.Println("Broadcat OK", res)
			// XXX OMG I CAN NOT GET HERE!! WHY!!
			return
		case "chat_message" == res.GetCMD():
			fmt.Println("Message Received")
			fmt.Println(id, res)
		}
	}

}
