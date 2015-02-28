package handler

import (
	"fmt"
	"github.com/tomyhero/go-tcp_server/authorizer"
	"github.com/tomyhero/go-tcp_server/context"
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

func (h *EchoHandler) HookInitialize(database map[string]interface{}) {
}

func (h *EchoHandler) HookDestroy(database map[string]interface{}) {
}

func (h *EchoHandler) HookBeforeExecute(c *context.Context) {
	fmt.Println("Called BeforeExecuteHandler", c.Session, c.Database["echo"])
	session := c.Session
	_, ok := session["num"]
	if !ok {
		session["num"] = 0
	}
	session["num"] = session["num"].(int) + 1
}

func (h *EchoHandler) HookAfterExecute(c *context.Context) {
	fmt.Println("Called AfterExecuteHandler")
}

func (h *EchoHandler) ActionEcho(c *context.Context) (*context.Context, error) {
	c.Res.Body = c.Req.Body
	fmt.Println(c, "Echo Echo!", c.Session["num"].(int))

	for conn, _ := range c.Conns {
		err := c.CDataManager.Send(conn, c.Res.GetData())
		if err != nil {
			fmt.Println("Fail to send : %s", err)
			//glog.Warningf("Fail to send : %s", err)
		}
	}

	c.OnSendResponse = false

	return c, nil
}
