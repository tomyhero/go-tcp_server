package context

import (
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
	HookFinalize(g map[string]interface{}, myStore map[string]interface{})
	HookBeforeExecute(c *Context)
	HookAfterExecute(c *Context)
}

type IAuthorizer interface {
	Login(c *Context) bool
	Auth(c *Context) bool
}

type Context struct {
	Req     *CData
	Res     *CData
	Stash   map[string]interface{}
	GStore  map[string]interface{}
	myStore map[string]interface{}
}

type CData struct {
	Header map[string]interface{}
	Body   map[string]interface{}
}

func NewContext(gstore map[string]interface{}, data map[string]interface{}) (*Context, error) {
	req, err := CreateReq(data)
	if err != nil {
		return nil, err
	}
	return &Context{GStore: gstore, Req: req, Res: CreateRes(), Stash: map[string]interface{}{}}, nil
}

func (c *Context) SetupMyStore() {
	prefix := strings.Split(c.Req.GetCMD(), "_")[0]
	c.myStore = c.GStore[prefix].(map[string]interface{})
}
func (c *Context) MyStore() map[string]interface{} {
	return c.myStore
}

func CreateRes() *CData {
	return &CData{Header: map[string]interface{}{"STATUS": STATUS_FORBIDDEN}, Body: map[string]interface{}{}}
}

func CreateReq(data map[string]interface{}) (*CData, error) {
	return &CData{Header: data["H"].(map[string]interface{}), Body: data["B"].(map[string]interface{})}, nil
}

func (r *CData) GetCMD() string {
	return r.Header["CMD"].(string)
}

func (r *CData) GetData() map[string]interface{} {
	data := map[string]interface{}{}
	data["H"] = r.Header
	data["B"] = r.Body
	return data
}
