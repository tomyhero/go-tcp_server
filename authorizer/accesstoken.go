package authorizer

import (
	"github.com/golang/glog"
	"github.com/tomyhero/go-tcp_server/context"
	"github.com/tomyhero/go-tcp_server/util"
)

type AccessToken struct {
}

func (a AccessToken) Reconnect(c *context.Context) bool {
	// TODO
	return true
}
func (a AccessToken) Login(c *context.Context) bool {
	session := c.Session
	_, has := session["access_token"]
	if !has {
		session["access_token"] = map[string]interface{}{}
	}
	tokenStore := session["access_token"].(map[string]interface{})
	uid, err := util.GenUUID()
	if err != nil {
		glog.Warningf("Fail to get UID %s", err)
		return false
	}
	tokenStore[uid] = map[string]interface{}{}
	c.Res.Body["AUTH_ACCESS_TOKEN"] = uid

	myConn := c.Conns[c.Conn].(map[string]interface{})
	myConn["AUTH_ACCESS_TOKEN"] = uid

	return true
}

func (a AccessToken) Auth(c *context.Context) bool {
	accessToken, has := c.Req.Header["AUTH_ACCESS_TOKEN"]
	if !has {
		return false
	}
	if accessToken == c.Conns[c.Conn].(map[string]interface{})["AUTH_ACCESS_TOKEN"] {
		return true
	}
	return false
}
