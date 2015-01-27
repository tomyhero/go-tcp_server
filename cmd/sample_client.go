package main

import (
	"bufio"
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

	in := map[string]interface{}{"h": map[string]interface{}{"cmd": "echo_Echo"}, "b": map[string]interface{}{"text": "Hello World\n"}}
	serialize := serializer.MessagePack{}
	buf, err := serialize.Serialize(in)
	conn.Write(buf.Bytes())
	res, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(serialize.Deserialize(bytes.NewBufferString(res)))

}
