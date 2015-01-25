package sample

import (
	"fmt"
	"github.com/tomyhero/ore_server/context"
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

func (h *EchoHandler) ActionEcho(c *context.Context) (*context.Context, error) {
	//c.Res.Body = c.Req.Body
	return c, nil
}

func (h *EchoHandler) ActionTest() {
	fmt.Println("Hello World")
}

func (h *EchoHandler) Prefix() string {
	return "echo"
}
