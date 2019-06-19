package clear

import (
	"encoding/json"
)

type (
	Protocol struct{}
)

func (p *Protocol) Encode(obj interface{}) ([]byte, error) {
	return json.MarshalIndent(obj, "", "\t")
}

func (p *Protocol) Decode(data []byte, objPtr interface{}) error {
	return json.Unmarshal(data, objPtr)
}

func (p *Protocol) Type() string {
	return "application/json"
}
