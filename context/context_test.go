package context

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewContext(t *testing.T) {
	data := map[string]interface{}{"h": map[string]interface{}{"cmd": "prefix_Echo"}, "b": map[string]interface{}{"text": "Hello World\n"}}
	c, err := NewContext(data)
	fmt.Println(c)
	assert.Nil(t, err)
	assert.Equal(t, "prefix_Echo", c.Req.Header["cmd"])
}

func TestNewRequest(t *testing.T) {
	data := map[string]interface{}{"h": map[string]interface{}{"cmd": "prefix_Echo"}, "b": map[string]interface{}{"text": "Hello World\n", "id": []int{1, 2, 3}}}
	req, err := NewRequest(data)
	assert.Nil(t, err)
	assert.Equal(t, "prefix_Echo", req.Header["cmd"])

	s := reflect.ValueOf(req.Body["id"])
	ids := make([]int, s.Len())
	for i := 0; i < s.Len(); i++ {
		ss := reflect.ValueOf(s.Index(i).Interface())
		ids[i] = int(ss.Int())
	}
	assert.Equal(t, []int{1, 2, 3}, ids)
}
