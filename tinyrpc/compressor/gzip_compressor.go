package compressor

import (
	"bytes"
	"compress/gzip"
	"io"
	"log"
)

type GzipCompressor struct{}

func (_ GzipCompressor) Compress(data []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	w := gzip.NewWriter(buf)
	defer func(w *gzip.Writer) {
		err := w.Close()
		if err != nil {
			log.Fatal(err)
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
	return buf.Bytes(), err
}

func (_ GzipCompressor) Decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer func(r *gzip.Reader) {
		err := r.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(r)
	data, err = io.ReadAll(r)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	return data, nil
}
