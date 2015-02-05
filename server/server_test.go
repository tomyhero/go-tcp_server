package server

import (
	//"bytes"
	"fmt"
	"github.com/tomyhero/ore_server/context"
	"github.com/tomyhero/ore_server/example/handler"
	"github.com/tomyhero/ore_server/serializer"
	"net"
	"testing"
	"time"
)

func emptyPort() (int, error) {
	l, err := net.Listen("tcp", ":0")
	defer l.Close()

	if err != nil {
		return 0, fmt.Errorf("Fail to listen empty port")
	}

	addr := l.Addr()
	port := addr.(*net.TCPAddr).Port
	return port, nil
}

func TestClient(t *testing.T) {

	port, err := emptyPort()
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}

	sv := Server{Port: port}
	handlers := make([]context.IHandler, 1)
	handlers[0] = handler.NewEchoHandler()
	sv.Setup(handlers)
	go sv.Run()

	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Println(err)
		t.Fail()
		return
	}
	defer conn.Close()
	doEcho(conn)
	sv.Shutdown()
}

func doEcho(conn net.Conn) {
	in := map[string]interface{}{"H": map[string]interface{}{"CMD": "echo_Echo", "AUTH_PLAIN_PASSWORD": "1111"}, "B": map[string]interface{}{"text": "Hello World\n"}}
	serialize := serializer.MessagePack{}
	buf, err := serialize.Serialize(in)
	conn.Write(buf.Bytes())
	res := make([]byte, 2024)
	conn.Read(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println(serialize.Deserialize(bytes.NewBuffer(res)))
}
