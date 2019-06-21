package camo

func Encode(data []byte) ([]byte, error) {
	data, err := compress(data)
	if err != nil {
		return nil, err
	}
	return encrypt(data)
}

func Decode(data []byte) ([]byte, error) {
	data, err := decrypt(data)
	if err != nil {
		return nil, err
	}
	return decompress(data)
}
