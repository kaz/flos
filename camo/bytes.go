package camo

import (
	"bytes"
	"io/ioutil"
)

func Encode(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	ew, err := encrypt(buf)
	if err != nil {
		return nil, err
	}

	cw := compress(ew)
	if _, err := cw.Write(data); err != nil {
		return nil, err
	}
	if err := cw.Flush(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Decode(data []byte) ([]byte, error) {
	dr, err := decrypt(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(decompress(dr))
}
