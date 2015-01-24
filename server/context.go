package server

import (
	"bytes"
	"github.com/tomyhero/ore_server/util"
)

type Context struct {
	Req   *Request
	Res   *Response
	Stash map[string]interface{}
}

type Request struct {
	Header map[string]interface{}
	Body   map[string]interface{}
}

type Response struct {
	Header map[string]interface{}
	Body   map[string]interface{}
}

func NewContext(buf *bytes.Buffer) (*Context, error) {
	req, err := NewRequest(buf)
	if err != nil {
		return nil, err
	}
	return &Context{Req: req}, nil
}

func NewRequest(buf *bytes.Buffer) (*Request, error) {
	data, err := util.UnpackMP(buf)
	if err != nil {
		return nil, err
	}
	return &Request{Header: data["h"].(map[string]interface{}), Body: data["b"].(map[string]interface{})}, nil
}
