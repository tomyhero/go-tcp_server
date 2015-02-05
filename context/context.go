package context

import (
	"net"
	"strings"
)

const (
	STATUS_NOT_PREPARE       = 0
	STATUS_OK                = 1
	STATUS_FORBIDDEN         = 2
	STATUS_NOT_OK            = 3
	STATUS_COMMAND_NOT_FOUND = 4
)

type IHandler interface {
	Prefix() string
	GetAuthorizer() IAuthorizer
	HookInitialize(g map[string]interface{}, myStore map[string]interface{})
	HookDestroy(g map[string]interface{}, myStore map[string]interface{})
	HookBeforeExecute(c *Context)
	HookAfterExecute(c *Context)
}

type IAuthorizer interface {
	Login(c *Context) bool
	Auth(c *Context) bool
}

type Context struct {
	Req          *CData
	Res          *CData
	Conn         net.Conn
	CDataManager *CDataManager
	GStore       map[string]interface{}
	myStore      map[string]interface{}
}

func NewContext(conn net.Conn, gstore map[string]interface{}, data map[string]interface{}) (*Context, error) {
	req, err := CreateReq(data)
	if err != nil {
		return nil, err
	}
	context := &Context{
		Conn:         conn,
		GStore:       gstore,
		Req:          req,
		Res:          CreateRes(req.GetCMD()),
		CDataManager: &CDataManager{},
	}
	return context, nil
}

func (c *Context) SetupMyStore() {
	prefix := strings.Split(c.Req.GetCMD(), "_")[0]
	c.myStore = c.GStore[prefix].(map[string]interface{})
}
func (c *Context) MyStore() map[string]interface{} {
	return c.myStore
}

func CreateRes(reqCmd string) *CData {
	cmd := reqCmd + "_res"
	return &CData{Header: map[string]interface{}{"STATUS": STATUS_FORBIDDEN, "CMD": cmd}, Body: map[string]interface{}{}}
}

func CreateReq(data map[string]interface{}) (*CData, error) {
	return &CData{Header: data["H"].(map[string]interface{}), Body: data["B"].(map[string]interface{})}, nil
}
