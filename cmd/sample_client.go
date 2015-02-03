package main

import (
	"bytes"
	"fmt"
	"github.com/tomyhero/ore_server/serializer"
	"net"
)

func main() {

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	doLogin(conn)
	doEcho(conn)

}

func doLogin(conn net.Conn) {
	in := map[string]interface{}{"h": map[string]interface{}{"cmd": "echo_login", "plain_password": "1111"}, "b": map[string]interface{}{}}
	serialize := serializer.MessagePack{}
	buf, err := serialize.Serialize(in)
	conn.Write(buf.Bytes())
	res := make([]byte, 2024)
	conn.Read(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(serialize.Deserialize(bytes.NewBuffer(res)))
}

func doEcho(conn net.Conn) {
	in := map[string]interface{}{"h": map[string]interface{}{"cmd": "echo_Echo"}, "b": map[string]interface{}{"text": "Hello World\n"}}
	serialize := serializer.MessagePack{}
	buf, err := serialize.Serialize(in)
	conn.Write(buf.Bytes())
	res := make([]byte, 2024)
	conn.Read(res)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(serialize.Deserialize(bytes.NewBuffer(res)))
}
