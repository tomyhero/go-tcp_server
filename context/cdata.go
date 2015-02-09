package context

import (
	//"bytes"
	"errors"
	"fmt"
	"github.com/tomyhero/ore_server/serializer"
	"github.com/ugorji/go/codec"
	"net"
	"reflect"
)

const (
	SERIALIZOR_TYPE_MESSAGE_PACK = 0
	SERIALIZOR_TYPE_JSON         = 1
)

type CDataManager struct {
	SerializorType int
}

type CData struct {
	Header map[string]interface{}
	Body   map[string]interface{}
}

func (r *CData) GetCMD() string {
	return r.Header["CMD"].(string)
}

func (r *CData) GetData() map[string]interface{} {
	data := map[string]interface{}{}
	data["H"] = r.Header
	data["B"] = r.Body
	return data
}

func (c CDataManager) Receive(conn net.Conn) (data map[string]interface{}, err error) {
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

	mh := codec.MsgpackHandle{}
	mh.MapType = reflect.TypeOf(map[string]interface{}{})
	mh.RawToString = true
	dec := codec.NewDecoder(conn, &mh)
	err = dec.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c CDataManager) Send(conn net.Conn, data map[string]interface{}) error {

	if c.SerializorType == SERIALIZOR_TYPE_MESSAGE_PACK {
		serializer := serializer.MessagePack{}
		buf, err := serializer.Serialize(data)
		if err != nil {
			return err
		}
		fmt.Println("SEND", conn, data["H"].(map[string]interface{})["CMD"])
		_, err = conn.Write(buf.Bytes())
		if err != nil {
			return err
		}
	} else if c.SerializorType == SERIALIZOR_TYPE_JSON {
		serializer := serializer.JSON{}
		buf, err := serializer.Serialize(data)
		if err != nil {
			return err
		}
		_, err = conn.Write(buf.Bytes())
	}

	return nil
}
