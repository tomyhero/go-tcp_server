package cdata

import (
	"bytes"
	"fmt"
	"gopkg.in/vmihailenco/msgpack.v2"
)

type MessagePack struct{}

func (m MessagePack) Encode(in map[string]interface{}) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	enc := msgpack.NewEncoder(buf)
	enc.Encode(in)
	return buf, nil
}

func (m MessagePack) Decode(buf *bytes.Buffer) (map[string]interface{}, error) {
	dec := msgpack.NewDecoder(buf)
	dec.DecodeMapFunc = func(d *msgpack.Decoder) (interface{}, error) {
		n, err := d.DecodeMapLen()
		if err != nil {
			return nil, err
		}
		// XXX all keys must be string.
		m := make(map[string]interface{}, n)
		for i := 0; i < n; i++ {
			mk, err := d.DecodeString()
			if err != nil {
				return nil, err
			}
			mv, err := d.DecodeInterface()
			if err != nil {
				return nil, err
			}
			m[mk] = mv
		}
		return m, nil
	}

	out, err := dec.DecodeInterface()
	return out.(map[string]interface{}), err
}

// maybe replace to more seriouls one.
func ValidateData(data map[string]interface{}) error {

	h, ok := data["h"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("header section not found")
	}

	_, ok = data["b"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("b section not found")
	}

	cmd, ok := h["cmd"]
	if !ok {
		return fmt.Errorf("header.cmd section not found", cmd)
	}
	/*
		if reflect.TypeOf(cmd) != string {
				matched, err := regexp.MatchString("^[a-zA-Z0-9_]*$", cmd)
				return fmt.Errorf("header.cmd must be string")
		}
	*/

	return nil
}
