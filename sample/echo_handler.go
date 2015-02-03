package sample

import (
	"fmt"
	"github.com/tomyhero/ore_server/authorizer"
	"github.com/tomyhero/ore_server/context"
)

type EchoHandler struct {
	HookHandler interface{}
	Authorizer  context.IAuthorizer
}

func (h *EchoHandler) Prefix() string {
	return "echo"
}
func (h *EchoHandler) AuthorizerHandler() context.IAuthorizer {
	return h.Authorizer
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
	c.Res.Body = c.Req.Body
	fmt.Println(c, "Echo Echo!")
	return c, nil
}
