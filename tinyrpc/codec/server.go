package codec

import (
	"bufio"
	"github.com/braver-braver/tinyrpc/compressor"
	"github.com/braver-braver/tinyrpc/header"
	"github.com/braver-braver/tinyrpc/serializer"
	"hash/crc32"
	"io"
	"net/rpc"
	"sync"
)

type reqCtx struct {
	requestID   uint64
	compareType compressor.CompressType
}

type serverCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	requestHeader header.RequestHeader
	serializer    serializer.Serializer
	mutex         sync.Mutex
	seq           uint64
	pending       map[uint64]*reqCtx
}

func NewServerCodec(conn io.ReadWriteCloser, serializer serializer.Serializer) rpc.ServerCodec {
	return &serverCodec{
		r:          bufio.NewReader(conn),
		w:          bufio.NewWriter(conn),
		c:          conn,
		serializer: serializer,
		pending:    make(map[uint64]*reqCtx),
	}
}

// ReadRequestHeader reads the rpc request header from io stream.
func (s *serverCodec) ReadRequestHeader(r *rpc.Request) error {
	s.requestHeader.ResetHeader()
	data, err := receiveFrame(s.r)
	if err != nil {
		return err
	}
	err = s.requestHeader.Unmarshall(data)
	if err != nil {
		return err
	}
	s.mutex.Lock()
	s.seq++
	s.pending[s.seq] = &reqCtx{
		s.requestHeader.ID,
		s.requestHeader.GetCompressType(),
	}
	r.ServiceMethod = s.requestHeader.Method
	r.Seq = s.seq
	s.mutex.Unlock()
	return nil
}

// ReadRequestBody reads the rpc request body from io stream.
func (s *serverCodec) ReadRequestBody(param any) error {
	if param == nil {
		if s.requestHeader.RequestLen != 0 {
			if err := read(s.r, make([]byte, s.requestHeader.RequestLen)); err != nil {
				return err
			}
		}
		return nil
	}

	reqBody := make([]byte, s.requestHeader.RequestLen)
	err := read(s.r, reqBody)
	if err != nil {
		return err
	}
	if _, ok := compressor.Compressors[s.requestHeader.CompressType]; !ok {
		return NotFoundCompressorError
	}
	if s.requestHeader.Checksum != 0 && crc32.ChecksumIEEE(reqBody) != s.requestHeader.Checksum {
		return UnexpectedChecksumError
	}

	req, err := compressor.Compressors[s.requestHeader.CompressType].Decompress(reqBody)
	if err != nil {
		return err
	}
	return s.serializer.UnMarshal(req, param)
}

// WriteResponse Write the rpc response header and body to the io stream
func (s *serverCodec) WriteResponse(response *rpc.Response, param any) error {
	s.mutex.Lock()
	reqCtx, ok := s.pending[response.Seq]
	if !ok {
		s.mutex.Unlock()
		return InvalidSequenceError
	}
	delete(s.pending, response.Seq)
	s.mutex.Unlock()

	if response.Error != "" {
		param = nil
	}
	if _, ok := compressor.Compressors[reqCtx.compareType]; !ok {
		return NotFoundCompressorError
	}

	var respBody []byte
	var err error
	if param != nil {
		respBody, err = s.serializer.Marshal(param)
		if err != nil {
			return err
		}
	}

	compressedResponseBody, err := compressor.Compressors[reqCtx.compareType].Compress(respBody)
	if err != nil {
		return err
	}

	h := header.ResponsePool.Get().(*header.ResponseHeader)
	defer func() {
		h.ResetHeader()
		header.ResponsePool.Put(h)
	}()

	h.ID = reqCtx.requestID
	h.Error = response.Error
	h.Checksum = crc32.ChecksumIEEE(compressedResponseBody)
	h.CompressType = reqCtx.compareType
	h.ResponseLen = uint32(len(compressedResponseBody))

	if err = sendFrame(s.w, h.Marshall()); err != nil {
		return err
	}

	if err = write(s.w, compressedResponseBody); err != nil {
		return err
	}
	err = s.w.(*bufio.Writer).Flush()
	if err != nil {
		return err
	}
	return nil
}

func (s *serverCodec) Close() error {
	return s.c.Close()
}
