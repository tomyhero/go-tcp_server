package handler

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/tomyhero/go-tcp_server/authorizer"
	"github.com/tomyhero/go-tcp_server/context"
	//"time"
)

var authAccessToken = authorizer.AccessToken{}

// Authorizer
type ChatAuthorizer struct {
}

func (a ChatAuthorizer) Login(c *context.Context) bool {
	name, hasName := c.Res.Body["name"]
	if hasName == false {
		name = "Unknown"
	}
	c.MyStore()["name"] = name

	return authAccessToken.Login(c)
}
func (a ChatAuthorizer) Auth(c *context.Context) bool {
	return authAccessToken.Auth(c)
}

// Setup Section

type ChatHandler struct {
	Authorizer context.IAuthorizer
}

func (h *ChatHandler) Prefix() string {
	return "chat"
}

func (h *ChatHandler) GetAuthorizer() context.IAuthorizer {
	return h.Authorizer
}

func NewChatHandler() *ChatHandler {
	return &ChatHandler{Authorizer: ChatAuthorizer{}}
}

// HOOK Section

func (h *ChatHandler) HookInitialize(g map[string]interface{}, myStore map[string]interface{}) {
}
func (h *ChatHandler) HookDestroy(g map[string]interface{}, myStore map[string]interface{}) {
}

func (h *ChatHandler) HookBeforeExecute(c *context.Context) {
}

func (h *ChatHandler) HookAfterExecute(c *context.Context) {
}

// Action Section

func (h *ChatHandler) ActionBroadcast(c *context.Context) {

	fmt.Println("Name", c.MyStore()["name"])
	cdata := context.CData{
		Header: map[string]interface{}{"CMD": "chat_message"},
		Body:   map[string]interface{}{"from": c.MyStore()["name"], "message": c.Req.Body["message"]},
	}

	for conn, _ := range c.ConnStore {
		err := c.CDataManager.Send(conn, cdata.GetData())
		if err != nil {
			if c.Conn == conn {
				glog.Warningf("Fail to send myself : %s", err)
			} else {
				glog.Warningf("Fail to send : %s", err)
			}
		}
	}
}
