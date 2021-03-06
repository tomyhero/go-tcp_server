package authorizer

import (
	"github.com/tomyhero/go-tcp_server/context"
)

type PlainPassword struct {
	Password string
}

func (a PlainPassword) Reconnect(c *context.Context) bool {
	return true
}

func (a PlainPassword) Login(c *context.Context) bool {
	// do nothing
	return true
}

func (a PlainPassword) Auth(c *context.Context) bool {
	password := c.Req.Header["AUTH_PLAIN_PASSWORD"]
	if password == a.Password {
		return true
	} else {
		return false
	}
}
