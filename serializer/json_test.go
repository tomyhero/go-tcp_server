package serializer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSON(t *testing.T) {
	in := map[string]interface{}{"h": map[string]interface{}{"cmd": "prefix_Echo"}, "b": map[string]interface{}{"text": "Hello World\n"}}
	serialize := JSON{}
	buf, err := serialize.Serialize(in)
	assert.Nil(t, err)
	out, err := serialize.Deserialize(buf)
	assert.Nil(t, err)
	assert.Equal(t, "prefix_Echo", out["h"].(map[string]interface{})["cmd"])
}
