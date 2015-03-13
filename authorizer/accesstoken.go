package authorizer

import (
	"github.com/tomyhero/go-tcp_server/context"
)

type AccessToken struct {
}

func (a AccessToken) Reconnect(c *context.Context) bool {

	/*
		sessionID, has := c.Res.Body["SESSION_ID"]
		if !has {
			return false
		}
	*/

	// TODO

	return true
}
func (a AccessToken) Login(c *context.Context) bool {
	c.Res.Body["AUTH_ACCESS_TOKEN"] = c.SessionID()
	return true
}

func (a AccessToken) Auth(c *context.Context) bool {
	accessToken, has := c.Req.Header["AUTH_ACCESS_TOKEN"]
	if !has {
		return false
	}
	if accessToken == c.Conns[c.Conn].(map[string]interface{})["session_id"] {
		return true
	}
	return false
}
