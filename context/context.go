package context

const (
	STATUS_NOT_PREPARE       = 0
	STATUS_OK                = 1
	STATUS_FORBIDDEN         = 2
	STATUS_NOT_OK            = 3
	STATUS_COMMAND_NOT_FOUND = 4
)

type IHandler interface {
	Prefix() string
	AuthorizerHandler() IAuthorizer
	BeforeExecuteHandler(c *Context)
	AfterExecuteHandler(c *Context)
}

type IAuthorizer interface {
	Login(c *Context) bool
	Auth(c *Context) bool
}

type Context struct {
	Req   *Request
	Res   *Response
	Stash map[string]interface{}
}

type Request struct {
	Header map[string]interface{}
	Body   map[string]interface{}
}

func (r *Request) GetCMD() string {
	return r.Header["CMD"].(string)
}

type Response struct {
	Header map[string]interface{}
	Body   map[string]interface{}
}

func NewContext(data map[string]interface{}) (*Context, error) {
	req, err := NewRequest(data)
	if err != nil {
		return nil, err
	}
	return &Context{Req: req, Res: NewResponse(), Stash: map[string]interface{}{}}, nil
}

func NewResponse() *Response {
	return &Response{Header: map[string]interface{}{"STATUS": STATUS_FORBIDDEN}, Body: map[string]interface{}{}}
}
func (r *Response) GetData() map[string]interface{} {
	data := map[string]interface{}{}
	data["H"] = r.Header
	data["B"] = r.Body
	return data
}

func NewRequest(data map[string]interface{}) (*Request, error) {
	return &Request{Header: data["H"].(map[string]interface{}), Body: data["B"].(map[string]interface{})}, nil
}
