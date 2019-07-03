package camo

import (
	"bytes"
	"io"
	"io/ioutil"
)

func EncodeWriter(w io.Writer) (io.WriteCloser, error) {
	return encrypt(w)
}
func Encode(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, len(data)))

	encWriter, err := EncodeWriter(buf)
	if err != nil {
		return nil, err
	}
	if _, err := encWriter.Write(data); err != nil {
		return nil, err
	}
	if err := encWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func DecodeReader(r io.Reader) (io.Reader, error) {
	return decrypt(r)
}
func Decode(data []byte) ([]byte, error) {
	decReader, err := DecodeReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(decReader)
}
