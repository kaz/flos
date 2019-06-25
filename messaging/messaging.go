package messaging

import (
	"log"
	"os"

	"github.com/kaz/flos/messaging/clear"
	"github.com/kaz/flos/messaging/color"
)

type (
	Protocol interface {
		Encode(interface{}) ([]byte, error)
		Decode([]byte, interface{}) error
		Type() string
	}
)

var (
	logger = log.New(os.Stdout, "[messaging] ", log.Ltime)

	DefaultProtocol Protocol
)

func Init() {
	if os.Getenv("FLOS_PROTO") == "clear" {
		DefaultProtocol = &clear.Protocol{}
		logger.Println("Using clear protocol")
	} else {
		DefaultProtocol = &color.Protocol{}
		logger.Println("Using color protocol")
	}
}

func Encode(obj interface{}) ([]byte, error) {
	return DefaultProtocol.Encode(obj)
}
func Decode(data []byte, objPtr interface{}) error {
	return DefaultProtocol.Decode(data, objPtr)
}
func Type() string {
	return DefaultProtocol.Type()
}
