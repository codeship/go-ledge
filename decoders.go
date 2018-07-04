package ledge

import "bufio"

var (
	rpcDecoderInstance = &rpcDecoder{}
)

type rpcDecoder struct{}

func (r *rpcDecoder) Decode(bufReader *bufio.Reader) ([]byte, error) {
	return bufReader.ReadBytes(separator)
}
