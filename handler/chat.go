package handler

import (
	"github.com/golang/glog"
	"github.com/tomyhero/go-tcp_server/authorizer"
	"github.com/tomyhero/go-tcp_server/context"
	"net"
)

var authAccessToken = authorizer.AccessToken{}

// Authorizer
// 名前でログインして、アクセストークン認証を活用する(AUTH_ACCESS_TOKENを返却するので、それを次回からヘッダーにつけて認証をおこなう）
type ChatAuthorizer struct{}

func (a ChatAuthorizer) Login(c *context.Context) bool {
	name, hasName := c.Res.Body["name"]
	if hasName == false {
		name = "Unknown"
	}
	c.Session["name"] = name

	return authAccessToken.Login(c)
}
func (a ChatAuthorizer) Auth(c *context.Context) bool {
	return authAccessToken.Auth(c)
}

// チャットサーバ準備
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

// チャットサーバのHOOK。必要に応じて活用。

func (h *ChatHandler) HookInitialize(database map[string]interface{}, conns map[net.Conn]interface{}) {
}
func (h *ChatHandler) HookDestroy(database map[string]interface{}) {
}

func (h *ChatHandler) HookBeforeExecute(c *context.Context) {
}

func (h *ChatHandler) HookAfterExecute(c *context.Context) {
}
func (h *ChatHandler) HookDisconnect(conn net.Conn, database map[string]interface{}, conns map[net.Conn]interface{}) {

}

// チャットのロジック。

func (h *ChatHandler) ActionBroadcast(c *context.Context) {

	cdata := context.CData{
		Header: map[string]interface{}{"CMD": "chat_message"},
		Body:   map[string]interface{}{"from": c.Session["name"], "message": c.Req.Body["message"]},
	}

	for conn, _ := range c.Conns {
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
