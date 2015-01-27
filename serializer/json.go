package serializer

import (
	"bytes"
	"encoding/json"
)

type JSON struct {
}

func (s JSON) Serialize(in map[string]interface{}) (*bytes.Buffer, error) {
	b, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)
	return buf, nil
}

func (s JSON) Deserialize(buf *bytes.Buffer) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	err := json.Unmarshal(buf.Bytes(), &data)
	return data, err
}
