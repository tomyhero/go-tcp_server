package server

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tomyhero/ore_server/serializer"
	"net"
)

const (
	SERIALIZOR_TYPE_MESSAGE_PACK = 0
	SERIALIZOR_TYPE_JSON         = 1
)

type CData struct {
	SerializorType int
}

func (c CData) Receive(conn net.Conn) (data map[string]interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			switch x := e.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			data = nil
			//fmt.Println("receive error: ", e) // Prints "Whoops: boom!"
		}
	}()

	b := make([]byte, 1024) // XXX
	_, err = conn.Read(b)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)

	if c.SerializorType == SERIALIZOR_TYPE_MESSAGE_PACK {
		serializer := serializer.MessagePack{}
		data, err = serializer.Deserialize(buf)
		if err != nil {
			fmt.Println(err)
		}
	} else if c.SerializorType == SERIALIZOR_TYPE_JSON {
		serializer := serializer.JSON{}
		data, err = serializer.Deserialize(buf)
		if err != nil {
			fmt.Println(err)
		}
	}

	return data, err
}

func (c CData) Send(conn net.Conn, data map[string]interface{}) error {

	if c.SerializorType == SERIALIZOR_TYPE_MESSAGE_PACK {
		serializer := serializer.MessagePack{}
		buf, err := serializer.Serialize(data)
		if err != nil {
			return err
		}
		conn.Write(buf.Bytes())
	} else if c.SerializorType == SERIALIZOR_TYPE_JSON {
		serializer := serializer.JSON{}
		buf, err := serializer.Serialize(data)
		if err != nil {
			return err
		}
		conn.Write(buf.Bytes())
	}

	return nil
}
