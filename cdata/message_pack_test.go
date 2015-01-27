package cdata

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMessagePack(t *testing.T) {
	in := map[string]interface{}{"h": map[string]interface{}{"cmd": "prefix_Echo"}, "b": map[string]interface{}{"text": "Hello World\n"}}
	cdata := MessagePack{}
	buf, err := cdata.Encode(in)
	assert.Nil(t, err)
	out, err := cdata.Decode(buf)
	assert.Nil(t, err)
	assert.Equal(t, "prefix_Echo", out["h"].(map[string]interface{})["cmd"])
}
