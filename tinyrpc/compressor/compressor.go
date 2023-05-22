package compressor

type CompressType uint16

const (
	Raw CompressType = iota
	Gzip
	Snappy
	Zlib
)

type Compressor interface {
	Compress([]byte) ([]byte, error)
	Decompress([]byte) ([]byte, error)
}

var Compressors = map[CompressType]Compressor{
	Raw:    RawCompressor{},
	Gzip:   GzipCompressor{},
	Snappy: SnappyCompressor{},
	Zlib:   ZlibCompressor{},
}
