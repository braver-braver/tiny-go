package tinyrpc

import (
	"github.com/braver-braver/tinyrpc/codec"
	"github.com/braver-braver/tinyrpc/compressor"
	"github.com/braver-braver/tinyrpc/serializer"
	"log"
	"net"
	"net/rpc"
)

// Option provides options for rpc
type Option func(o *options)

type options struct {
	compressType compressor.CompressType
	serializer   serializer.Serializer
}

type Server struct {
	*rpc.Server
	serializer.Serializer
}

func NewServer(opts ...Option) *Server {
	options := options{
		serializer: serializer.Proto,
	}
	for _, opt := range opts {
		opt(&options)
	}

	return &Server{&rpc.Server{}, options.serializer}
}

// Register register rpc function
func (s *Server) Register(rcvr interface{}) error {
	return s.Server.Register(rcvr)
}

// RegisterName register the rpc function with the specified name
func (s *Server) RegisterName(name string, rcvr interface{}) error {
	return s.Server.RegisterName(name, rcvr)
}

func (s *Server) Serve(lis net.Listener) error {
	log.Printf("tinyrpc server listening on %s", lis.Addr().String())

	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}
		go s.Server.ServeCodec(codec.NewServerCodec(conn, s.Serializer))
	}
}
