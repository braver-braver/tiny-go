package header

import (
	"encoding/binary"
	"errors"
	"github.com/braver-braver/tinyrpc/compressor"
	"sync"
)

const (
	MaxHeaderSize = 2 + 10 + 10 + 10 + 4
	Uint32Size    = 4
	Uint16Size    = 2
)

var UnmarshalError = errors.New("an error occurred in Unmarshal")

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
	idx += Uint16Size
	idx += writeString(header[idx:], r.Method)
	idx += binary.PutUvarint(header[idx:], r.ID)
	idx += binary.PutUvarint(header[idx:], uint64(r.RequestLen))
	binary.LittleEndian.PutUint32(header[idx:], r.Checksum)
	idx += Uint32Size
	return header[:idx]
}

// Unmarshall decode byte slice into RequestHeader structure
func (r *RequestHeader) Unmarshall(data []byte) (err error) {
	r.Lock()
	defer r.Unlock()
	if len(data) == 0 {
		return UnmarshalError
	}
	defer func() {
		if r := recover(); r != nil {
			err = UnmarshalError
		}
	}()
	idx, size := 0, 0
	r.CompressType = compressor.CompressType(binary.LittleEndian.Uint16(data[idx:]))
	idx += Uint16Size

	r.Method, size = readString(data[idx:])
	idx += size

	r.ID, size = binary.Uvarint(data[idx:])
	idx += size

	requestLen, size := binary.Uvarint(data[idx:])
	r.RequestLen = uint32(requestLen)
	idx += size

	r.Checksum = binary.LittleEndian.Uint32(data[idx:])
	return
}

func (r *RequestHeader) GetCompressType() compressor.CompressType {
	r.RLock()
	defer r.RUnlock()
	return r.CompressType
}

func (r *RequestHeader) ResetHeader() {
	r.Lock()
	defer r.Unlock()
	r.ID = 0
	r.CompressType = 0
	r.Method = ""
	r.RequestLen = 0
	r.Checksum = 0
}

// ResponseHeader request header structure looks like:
// +--------------+---------+----------------+-------------+----------+
// | CompressType |    ID   |      Error     | ResponseLen | Checksum |
// +--------------+---------+----------------+-------------+----------+
// |    uint16    | uvarint | uvarint+string |    uvarint  |  uint32  |
// +--------------+---------+----------------+-------------+----------+

type ResponseHeader struct {
	sync.RWMutex
	CompressType compressor.CompressType
	ID           uint64
	Error        string
	ResponseLen  uint32
	Checksum     uint32
}

func (r *ResponseHeader) Marshall() []byte {
	r.RLock()
	defer r.RUnlock()
	idx := 0
	header := make([]byte, MaxHeaderSize+len(r.Error))
	binary.LittleEndian.PutUint16(header[idx:], uint16(r.CompressType))
	idx += Uint16Size

	idx += binary.PutUvarint(header[idx:], r.ID)

	idx += writeString(header[idx:], r.Error)

	idx += binary.PutUvarint(header[idx:], uint64(r.ResponseLen))

	binary.LittleEndian.PutUint32(header[idx:], r.Checksum)
	idx += Uint32Size

	return header[:idx]
}

func (r *ResponseHeader) Unmarshall(data []byte) (err error) {
	r.Lock()
	defer r.Unlock()
	if len(data) == 0 {
		return UnmarshalError
	}
	defer func() {
		if r := recover(); r != nil {
			err = UnmarshalError
		}
	}()
	idx, size := 0, 0
	r.CompressType = compressor.CompressType(binary.LittleEndian.Uint16(data[idx:]))
	idx += Uint16Size

	r.ID, size = binary.Uvarint(data[idx:])
	idx += size

	r.Error, size = readString(data[idx:])
	idx += size

	responseLen, size := binary.Uvarint(data[idx:])
	r.ResponseLen = uint32(responseLen)
	idx += size

	r.Checksum = binary.LittleEndian.Uint32(data[idx:])
	return
}

// GetCompressType get compress type
func (r *ResponseHeader) GetCompressType() compressor.CompressType {
	r.RLock()
	defer r.RUnlock()
	return r.CompressType
}

// ResetHeader reset response header
func (r *ResponseHeader) ResetHeader() {
	r.Lock()
	defer r.Unlock()
	r.Error = ""
	r.ID = 0
	r.CompressType = 0
	r.Checksum = 0
	r.ResponseLen = 0
}

func writeString(data []byte, str string) int {
	idx := 0
	idx += binary.PutUvarint(data[idx:], uint64(len(str)))
	// that equals copy(data[idx:], str)
	// idx += len(str)
	idx += copy(data[idx:], str)
	return idx
}

func readString(data []byte) (string, int) {
	idx := 0
	length, size := binary.Uvarint(data[idx:])
	idx += size
	str := string(data[idx : idx+int(length)])
	return str, idx + len(str)
}
