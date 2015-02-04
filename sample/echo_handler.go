package sample

import (
	"fmt"
	"github.com/tomyhero/ore_server/authorizer"
	"github.com/tomyhero/ore_server/context"
)

type EchoHandler struct {
	Authorizer context.IAuthorizer
}

func (h *EchoHandler) Prefix() string {
	return "echo"
}
func (h *EchoHandler) GetAuthorizer() context.IAuthorizer {
	return h.Authorizer
}

func NewEchoHandler() *EchoHandler {
	return &EchoHandler{Authorizer: authorizer.PlainPassword{Password: "1111"}}
}

func (h *EchoHandler) HookInitialize(g map[string]interface{}, myStore map[string]interface{}) {
	myStore["num"] = 0
}
func (h *EchoHandler) HookFinalize(g map[string]interface{}, myStore map[string]interface{}) {

}

func (h *EchoHandler) HookBeforeExecute(c *context.Context) {
	fmt.Println("Called BeforeExecuteHandler", c.MyStore, c.GStore["echo"])
	myStore := c.MyStore()
	myStore["num"] = myStore["num"].(int) + 1
}

func (h *EchoHandler) HookAfterExecute(c *context.Context) {
	fmt.Println("Called AfterExecuteHandler")
}

func (h *EchoHandler) ActionEcho(c *context.Context) (*context.Context, error) {
	c.Res.Body = c.Req.Body
	fmt.Println(c, "Echo Echo!", c.MyStore()["num"].(int))
	return c, nil
}
