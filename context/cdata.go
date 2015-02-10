package context

import (
	"errors"
	"github.com/ugorji/go/codec"
	"net"
	"reflect"
)

type CDataManager struct {
	CodecHandle codec.Handle
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

func (c *CDataManager) codecHandle() codec.Handle {
	if c.CodecHandle == nil {
		var h = new(codec.MsgpackHandle)
		h.MapType = reflect.TypeOf(map[string]interface{}{})
		h.RawToString = true
		c.CodecHandle = h
	}

	return c.CodecHandle
}

func (c *CDataManager) Receive(conn net.Conn) (data map[string]interface{}, err error) {
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

	dec := codec.NewDecoder(conn, c.codecHandle())
	err = dec.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *CDataManager) Send(conn net.Conn, data map[string]interface{}) error {
	enc := codec.NewEncoder(conn, c.codecHandle())
	err := enc.Encode(data)
	return err
}
