package authorizer

import (
	"github.com/golang/glog"
	"github.com/tomyhero/go-tcp_server/context"
	"github.com/tomyhero/go-tcp_server/util"
)

type AccessToken struct {
}

func (a AccessToken) Login(c *context.Context) bool {
	myStore := c.MyStore()
	_, has := myStore["access_token"]
	if !has {
		myStore["access_token"] = map[string]interface{}{}
	}
	tokenStore := myStore["access_token"].(map[string]interface{})
	uid, err := util.GenUUID()
	if err != nil {
		glog.Warningf("Fail to get UID %s", err)
		return false
	}
	tokenStore[uid] = map[string]interface{}{}
	c.Res.Body["AUTH_ACCESS_TOKEN"] = uid

	myConnStore := c.ConnStore[c.Conn].(map[string]interface{})
	myConnStore["AUTH_ACCESS_TOKEN"] = uid

	return true
}

func (a AccessToken) Auth(c *context.Context) bool {
	accessToken, has := c.Req.Header["AUTH_ACCESS_TOKEN"]
	if !has {
		return false
	}
	if accessToken == c.ConnStore[c.Conn].(map[string]interface{})["AUTH_ACCESS_TOKEN"] {
		return true
	}
	return false
}
