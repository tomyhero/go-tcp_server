package context

type IHandler interface {
	Prefix() string
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
	return r.Header["cmd"].(string)
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
	return &Response{Header: map[string]interface{}{}, Body: map[string]interface{}{}}
}
func (r *Response) GetData() map[string]interface{} {
	data := map[string]interface{}{}
	data["h"] = r.Header
	data["b"] = r.Body
	return data
}

func NewRequest(data map[string]interface{}) (*Request, error) {
	return &Request{Header: data["h"].(map[string]interface{}), Body: data["b"].(map[string]interface{})}, nil
}
