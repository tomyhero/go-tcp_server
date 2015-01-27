package sample

import (
	"fmt"
	"github.com/tomyhero/ore_server/authorizer"
	"github.com/tomyhero/ore_server/context"
)

type EchoHandler struct {
	HookHandler interface{}
	Authorizer  authorizer.IAuthorizer
}

func (h *EchoHandler) Prefix() string {
	return "echo"
}

func NewEchoHandler() *EchoHandler {
	//&EchoHandler{HookHandler: HookHandler{}, AuthrizeHander: auth.NillAutuer{}}
	return &EchoHandler{Authorizer: authorizer.PlainPassword{Password: "1111"}}
}

type HookHandler struct{}

func (h *HookHandler) Initialize() {
}

func (h *HookHandler) Finalize() {
}

func (h *EchoHandler) ActionEcho(c *context.Context) (*context.Context, error) {
	//c.Res.Body = c.Req.Body
	fmt.Println("Echo Echo!")
	return c, nil
}
