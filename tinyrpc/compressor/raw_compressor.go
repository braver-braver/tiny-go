package compressor

type RawCompressor struct{}

func (_ RawCompressor) Compress(data []byte) ([]byte, error) {
	return data, nil
}

func (_ RawCompressor) Decompress(data []byte) ([]byte, error) {
	return data, nil
}
