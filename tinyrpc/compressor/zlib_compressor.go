package compressor

import (
	"bytes"
	"compress/zlib"
	"io"
)

type ZlibCompressor struct {
}

func (_ ZlibCompressor) Compress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	w := zlib.NewWriter(buf)
	defer func(w *zlib.Writer) {
		err := w.Close()
		if err != nil {

		}
	}(w)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (_ ZlibCompressor) Decompress(data []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {

		}
	}(reader)
	data, err = io.ReadAll(reader)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	return data, nil
}
