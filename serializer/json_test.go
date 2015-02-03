package serializer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSON(t *testing.T) {
	in := map[string]interface{}{"H": map[string]interface{}{"CMD": "prefix_Echo"}, "B": map[string]interface{}{"text": "Hello World\n"}}
	serialize := JSON{}
	buf, err := serialize.Serialize(in)
	assert.Nil(t, err)
	out, err := serialize.Deserialize(buf)
	assert.Nil(t, err)
	assert.Equal(t, "prefix_Echo", out["H"].(map[string]interface{})["CMD"])
}
