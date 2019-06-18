package messaging

import (
	"github.com/kaz/flos/messaging/clear"
)

type (
	Protocol interface {
		Encode(interface{}) ([]byte, error)
		Decode([]byte, interface{}) error
		Type() string
	}
)

var (
	DefaultProtocol = &clear.Protocol{}
	// DefaultProtocol = &daphne.Protocol{}
)

func Encode(obj interface{}) ([]byte, error) {
	return DefaultProtocol.Encode(obj)
}
func Decode(data []byte, objPtr interface{}) error {
	return DefaultProtocol.Decode(data, objPtr)
}
func Type() string {
	return DefaultProtocol.Type()
}
