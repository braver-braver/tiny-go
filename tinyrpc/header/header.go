package header

import (
	"encoding/binary"
	"github.com/braver-braver/tinyrpc/compressor"
	"sync"
)

const (
	MaxHeaderSize = 2 + 10 + 10 + 10 + 4
	Uint32Size    = 4
	Uint16Size    = 2
)

// RequestHeader request header structure looks like:
// +--------------+----------------+----------+------------+----------+
// | CompressType |      Method    |    ID    | RequestLen | Checksum |
// +--------------+----------------+----------+------------+----------+
// |    uint16    | uvarint+string |  uvarint |   uvarint  |  uint32  |
// +--------------+----------------+----------+------------+----------+
// for uvarint64, the default num length is 10

type RequestHeader struct {
	sync.RWMutex
	CompressType compressor.CompressType
	Method       string
	ID           uint64
	RequestLen   uint32 // ?? 请思考这里为什么使用的是 uint32 来对应 请求头结构的 uvarint
	Checksum     uint32
}

func (r *RequestHeader) Marshall() []byte {
	r.RLock()
	defer r.RUnlock()
	idx := 0
	header := make([]byte, MaxHeaderSize+len(r.Method))
	binary.LittleEndian.PutUint16(header[idx:], uint16(r.CompressType))
	idx += 2
	idx += writeString(header[idx:], r.Method)
	return nil
}

func writeString(data []byte, str string) int {
	idx := 0
	idx += binary.PutUvarint(data[idx:], uint64(len(str)))
	idx += copy(data[idx:], str)
	return idx
}
