package authorizer

import (
	"github.com/tomyhero/go-tcp_server/context"
)

type AccessToken struct {
}

func (a AccessToken) Reconnect(c *context.Context) bool {

	sessionID, has := c.Req.Header["AUTH_ACCESS_TOKEN"]

	if !has {
		return false
	}

	c.SetSessionID(sessionID.(string))

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

	if accessToken == c.SessionID() {
		return true
	}
	return false
}
