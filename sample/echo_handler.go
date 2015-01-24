package sample

import (
	. "github.com/tomyhero/ore_server/server"
)

type EchoHandler struct {
	HookHandler    interface{}
	AuthrizeHander interface{}
}

func NewEchoHandler() *EchoHandler {
	//&EchoHandler{HookHandler: HookHandler{}, AuthrizeHander: auth.NillAutuer{}}
	return &EchoHandler{}
}

type HookHandler struct{}

func (h *HookHandler) Initialize() {
}

func (h *HookHandler) Finalize() {
}

func (h *EchoHandler) Echo(c *Context) (*Context, error) {
	//c.Res.Body = c.Req.Body
	return c, nil
}
