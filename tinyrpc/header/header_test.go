package header

import (
	"github.com/braver-braver/tinyrpc/compressor"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRequestHeader_Marshall(t *testing.T) {
	header := &RequestHeader{
		CompressType: 0,
		Method:       "Add",
		ID:           12455,
		RequestLen:   266,
		Checksum:     3845236589,
	}

	assert.Equal(t, []byte{0x0, 0x0, 0x3, 0x41, 0x64, 0x64,
		0xa7, 0x61, 0x8a, 0x2, 0x6d, 0xa7, 0x31, 0xe5}, header.Marshall())
}

func TestRequestHeader_Unmarshall(t *testing.T) {
	type expect struct {
		header *RequestHeader
		err    error
	}
	cases := []struct {
		name   string
		data   []byte
		expect expect
	}{
		{
			"test-1",
			[]byte{0x0, 0x0, 0x3, 0x41, 0x64, 0x64,
				0xa7, 0x61, 0x8a, 0x2, 0x6d, 0xa7, 0x31, 0xe5},
			expect{&RequestHeader{
				CompressType: 0,
				Method:       "Add",
				ID:           12455,
				RequestLen:   266,
				Checksum:     3845236589,
			}, nil},
		},
		{
			"test-2",
			nil,
			expect{
				&RequestHeader{},
				UnmarshalError,
			},
		},
		{
			"test-3",
			[]byte{0x0},
			expect{&RequestHeader{}, UnmarshalError},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := &RequestHeader{}
			err := h.Unmarshall(c.data)
			assert.Equal(t, true, reflect.DeepEqual(c.expect.header, h))
			assert.Equal(t, err, c.expect.err)
		})
	}
}

func TestRequestHeader_ResetHeader(t *testing.T) {
	header := &RequestHeader{
		CompressType: 0,
		Method:       "Add",
		ID:           12455,
		RequestLen:   266,
		Checksum:     3845236589,
	}
	header.ResetHeader()
	assert.Equal(t, true, reflect.DeepEqual(header, &RequestHeader{}))
}

func TestResponseHeader_Marshall(t *testing.T) {
	header := &ResponseHeader{
		CompressType: 0,
		Error:        "error",
		ID:           12455,
		ResponseLen:  266,
		Checksum:     3845236589,
	}
	assert.Equal(t, []byte{0x0, 0x0, 0xa7, 0x61, 0x5, 0x65, 0x72,
		0x72, 0x6f, 0x72, 0x8a, 0x2, 0x6d, 0xa7, 0x31, 0xe5}, header.Marshall())
}

func TestResponseHeader_Unmarshall(t *testing.T) {
	type expect struct {
		header *ResponseHeader
		err    error
	}
	cases := []struct {
		name   string
		data   []byte
		expect expect
	}{
		{
			"test-1",
			[]byte{
				0x0, 0x0, 0xa7, 0x61, 0x5, 0x65, 0x72,
				0x72, 0x6f, 0x72, 0x8a, 0x2, 0x6d, 0xa7, 0x31, 0xe5,
			},
			expect{
				&ResponseHeader{
					CompressType: 0,
					Error:        "error",
					ID:           12455,
					ResponseLen:  266,
					Checksum:     3845236589,
				}, nil,
			},
		},
		{
			"test-2",
			nil,
			expect{
				&ResponseHeader{},
				UnmarshalError,
			},
		},
		{
			"test-3",
			[]byte{0x0},
			expect{
				&ResponseHeader{},
				UnmarshalError,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			h := &ResponseHeader{}
			err := h.Unmarshall(c.data)
			assert.Equal(t, true, reflect.DeepEqual(c.expect.header, h))
			assert.Equal(t, err, c.expect.err)
		})
	}
}

func TestResponseHeader_ResetHeader(t *testing.T) {
	header := &ResponseHeader{
		CompressType: 0,
		Error:        "error",
		ID:           12455,
		ResponseLen:  266,
		Checksum:     3845236589,
	}
	header.ResetHeader()
	assert.Equal(t, true, reflect.DeepEqual(header, &ResponseHeader{}))
}

func TestRequestHeader_GetCompressType(t *testing.T) {
	header := &RequestHeader{
		CompressType: 0,
		Method:       "Add",
		ID:           12455,
		RequestLen:   266,
		Checksum:     3845236589,
	}

	assert.Equal(t, true, reflect.DeepEqual(compressor.CompressType(0), header.GetCompressType()))
}

func TestResponseHeader_GetCompressType(t *testing.T) {
	header := &ResponseHeader{
		CompressType: 0,
		Error:        "error",
		ID:           12455,
		ResponseLen:  266,
		Checksum:     3845236589,
	}

	assert.Equal(t, true, reflect.DeepEqual(compressor.CompressType(0), header.GetCompressType()))
}
