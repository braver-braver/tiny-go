package serializer

import (
	"testing"
	"github.com/braver-braver/tinyrpc/mock/message"
)

type testStruct struct {

}
func TestProtoSerializer_Marshal(t *testing.T) {
	type expect struct {
		data []byte
		err error
	}
	cases := []struct{
		name string
		arg interface{}
		expect expect
	}{
		{
			"test-1",
			&pb.ArithmeticRequest{

			}
		},
	}

}

func TestProtoSerializer_UnMarshal(t *testing.T) {

}
