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

type clientCodec struct {
	r io.Reader
	w io.Writer
	c io.Closer

	compressor     compressor.CompressType
	serializer     serializer.Serializer
	responseHeader header.ResponseHeader
	mutex          sync.Mutex // protect pending map
	pending        map[uint64]string
}

func NewClientCodec(conn io.ReadWriteCloser, compressType compressor.CompressType, serializer serializer.Serializer) rpc.ClientCodec {
	return &clientCodec{
		r:          bufio.NewReader(conn),
		w:          bufio.NewWriter(conn),
		c:          conn,
		compressor: compressType,
		serializer: serializer,
		pending:    make(map[uint64]string),
	}
}

// WriteRequest writes a rpc requestHeader & its body  to io stream.
func (c *clientCodec) WriteRequest(r *rpc.Request, params interface{}) error {
	c.mutex.Lock()
	c.pending[r.Seq] = r.ServiceMethod
	c.mutex.Unlock()

	if _, ok := compressor.Compressors[c.compressor]; !ok {
		return NotFoundCompressorError
	}
	body, err := c.serializer.Marshal(params)
	if err != nil {
		return err
	}
	compressedBody, err := compressor.Compressors[c.compressor].Compress(body)
	if err != nil {
		return err
	}

	h := header.RequestPool.Get().(*header.RequestHeader)
	defer func() {
		h.ResetHeader()
		header.RequestPool.Put(h)
	}()
	h.ID = r.Seq
	h.Method = r.ServiceMethod
	h.RequestLen = uint32(len(compressedBody))
	h.CompressType = c.compressor
	h.Checksum = crc32.ChecksumIEEE(compressedBody)

	if err = sendFrame(c.w, h.Marshall()); err != nil {
		return err
	}
	// requestHeader 的 存在，已经标定了 请求体的长度，调用write 方法直接将 内容写入流即可。
	if err = write(c.w, compressedBody); err != nil {
		return err
	}

	err = c.w.(*bufio.Writer).Flush()
	if err != nil {
		return err
	}
	return nil
}

// ReadResponseHeader reads ResponseHeader from offered io stream
func (c *clientCodec) ReadResponseHeader(r *rpc.Response) error {
	c.responseHeader.ResetHeader()
	data, err := receiveFrame(c.r)
	if err != nil {
		return err
	}
	err = c.responseHeader.Unmarshall(data)
	if err != nil {
		return err
	}
	c.mutex.Lock()
	r.Seq = c.responseHeader.ID
	// TODD, _, ok := c.pending[r.Seq] ??
	r.ServiceMethod = c.pending[r.Seq]
	r.Error = c.responseHeader.Error
	delete(c.pending, r.Seq)
	c.mutex.Unlock()
	return nil
}

// ReadResponseBody reads rpc response body from offered io stream
func (c *clientCodec) ReadResponseBody(param any) error {
	if param == nil {
		if c.responseHeader.ResponseLen != 0 {
			if err := read(c.r, make([]byte, c.responseHeader.ResponseLen)); err != nil {
				return err
			}
		}
		return nil
	}

	responseBody := make([]byte, c.responseHeader.ResponseLen)
	err := read(c.r, responseBody)
	if err != nil {
		return err
	}
	if c.responseHeader.GetCompressType() != c.compressor {
		return CompressorTypeMismatchError
	}
	if c.responseHeader.Checksum != 0 && crc32.ChecksumIEEE(responseBody) != c.responseHeader.Checksum {
		return UnexpectedChecksumError
	}
	resp, err := compressor.Compressors[c.responseHeader.CompressType].Decompress(responseBody)
	if err != nil {
		return err
	}
	return c.serializer.UnMarshal(resp, param)
}

func (c *clientCodec) Close() error {
	return c.c.Close()
}
