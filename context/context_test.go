package context

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"reflect"
	"testing"
)

func TestNewContext(t *testing.T) {
	data := map[string]interface{}{"H": map[string]interface{}{"CMD": "prefix_Echo"}, "B": map[string]interface{}{"text": "Hello World\n"}}
	c, err := NewContext(nil, map[string]interface{}{}, data, map[net.Conn]interface{}{})
	fmt.Println(c)
	assert.Nil(t, err)
	assert.Equal(t, "prefix_Echo", c.Req.Header["CMD"])
}

func TestNewCData(t *testing.T) {
	data := map[string]interface{}{"H": map[string]interface{}{"CMD": "prefix_Echo"}, "B": map[string]interface{}{"text": "Hello World\n", "id": []int{1, 2, 3}}}
	req, err := CreateReq(data)
	assert.Nil(t, err)
	assert.Equal(t, "prefix_Echo", req.Header["CMD"])

	s := reflect.ValueOf(req.Body["id"])
	ids := make([]int, s.Len())
	for i := 0; i < s.Len(); i++ {
		ss := reflect.ValueOf(s.Index(i).Interface())
		ids[i] = int(ss.Int())
	}
	assert.Equal(t, []int{1, 2, 3}, ids)
}
