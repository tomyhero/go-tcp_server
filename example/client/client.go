package main

import (
	"fmt"
	"github.com/tomyhero/go-tcp_server/context"
	"net"
)

func main() {

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	//doLogin(conn)
	doEcho(conn)

}

func doLogin(conn net.Conn) {
	in := map[string]interface{}{"H": map[string]interface{}{"CMD": "echo_login", "plain_password": "1111"}, "B": map[string]interface{}{}}
	cm := &context.CDataManager{}
	cm.Send(conn, in)
	data, err := cm.Receive(conn)
	fmt.Println("doLogin", data, err)
}

func doEcho(conn net.Conn) {
	in := map[string]interface{}{"H": map[string]interface{}{"CMD": "echo_Echo", "AUTH_PLAIN_PASSWORD": "1111"}, "B": map[string]interface{}{"text": "Hello World\n"}}
	cm := &context.CDataManager{}
	cm.Send(conn, in)
	data, err := cm.Receive(conn)
	fmt.Println("doEcho", data, err)
}
