package codec

import (
	"encoding/binary"
	"io"
	"net"
)

// sendFrame 函数将会向IO 写入 uvarint类型的size，表示要发送数据的长度，随后将该字节slice 类型数据 data 写入 IO 流中。
func sendFrame(w io.Writer, data []byte) (err error) {
	var size [binary.MaxVarintLen64]byte

	if len(data) == 0 {
		n := binary.PutUvarint(size[:], uint64(0))
		if err = write(w, size[:n]); err != nil {
			return
		}
		return
	}

	n := binary.PutUvarint(size[:], uint64(len(data)))
	if err = write(w, size[:n]); err != nil {
		return err
	}

	if err = write(w, data); err != nil {
		return err
	}
	return
}

func receiveFrame(r io.Reader) (data []byte, err error) {
	size, err := binary.ReadUvarint(r.(io.ByteReader))
	if err != nil {
		return nil, err
	}
	if size != 0 {
		data = make([]byte, size)
		if err = read(r, data); err != nil {
			return nil, err
		}
	}
	return data, nil
}

// 注意， 实际实现中考虑传入的是 bufio 的 writer（或 reader）
// 注意，由于 codec 层会传入一个bufio类型的结构体，bufio类型实现了有缓冲的IO操作，
// 以便减少IO在用户态与内核态拷贝的次数。
func write(w io.Writer, data []byte) error {
	for index := 0; index < len(data); {
		n, err := w.Write(data[index:])
		if _, ok := err.(net.Error); !ok {
			return err
		}
		index += n
	}
	return nil
}

func read(r io.Reader, data []byte) (err error) {
	for index := 0; index < len(data); {
		n, err := r.Read(data[index:])
		if err != nil {
			if _, ok := err.(net.Error); !ok {
				return err
			}
		}
		index += n
	}
	return nil
}
