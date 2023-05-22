package compressor

import (
	"bytes"
	"github.com/klauspost/compress/snappy"
	"io"
)

type SnappyCompressor struct{}

func (_ SnappyCompressor) Compress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	w := snappy.NewBufferedWriter(buf)
	defer func(w *snappy.Writer) {
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

func (_ SnappyCompressor) Decompress(data []byte) ([]byte, error) {
	r := snappy.NewReader(bytes.NewBuffer(data))
	data, err := io.ReadAll(r)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	return data, nil
}
