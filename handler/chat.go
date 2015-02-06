package handler

import (
	"fmt"
	"github.com/tomyhero/ore_server/authorizer"
	"github.com/tomyhero/ore_server/context"
)

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
	return &ChatHandler{Authorizer: authorizer.AccessToken{}}
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

	cdata := context.CData{
		Header: map[string]interface{}{"CMD": "chat_message"},
		Body:   map[string]interface{}{"from": c.Req.Body["name"], "message": c.Req.Body["message"]},
	}

	for conn, _ := range c.ConnStore {
		err := c.CDataManager.Send(conn, cdata.GetData())
		if err != nil {
			fmt.Println(err)
		}
	}
}
