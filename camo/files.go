package camo

import (
	"fmt"
	"io/ioutil"
	"os"
)

func WriteFile(filename string, data []byte, perm os.FileMode) error {
	payload, err := Encode(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, append([]byte(CAMO_HEADER), payload...), perm)
}
func ReadFile(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if len(CAMO_HEADER) > len(data) {
		return nil, fmt.Errorf("file size too short")
	}

	return Decode(data[len(CAMO_HEADER):])
}
