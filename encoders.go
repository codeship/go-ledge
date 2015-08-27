package ledge

import (
	"errors"
	"io"
)

var (
	rpcEncoderInstance = &rpcEncoder{}

	separator                      = byte('\n')
	errShouldNotBeWritingZeroBytes = errors.New("ledge: should not be writing 0 bytes")
)

type rpcEncoder struct{}

func (r *rpcEncoder) Encode(writer io.Writer, data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}
	return writer.Write(append(data, separator))
}
