package ledge

import (
	"bufio"
	"fmt"
	"strconv"
)

var (
	rpcDecoderInstance = &rpcDecoder{}
)

type rpcDecoder struct{}

func (r *rpcDecoder) Decode(bufReader *bufio.Reader) ([]byte, error) {
	slice, err := bufReader.ReadSlice(separatorBytes[0])

	if err != nil {
		return nil, err
	}
	sizeBytesLenBuf := int64(slice[0])
	offset := int64(1)

	sizeBytes := slice[offset : offset+sizeBytesLenBuf]
	offset = offset + sizeBytesLenBuf
	size, err := strconv.ParseInt(string(sizeBytes), 10, 64)
	if err != nil {
		return nil, err
	}
	sizeInt := int(size)
	if int64(sizeInt) != size {
		return nil, fmt.Errorf("ledge: could not cast %d to an int", size)
	}

	data := slice[offset : offset+size]

	// docker logs do not flush without a newline, this is the hack
	if int64(len(data)) > size {
		return nil, fmt.Errorf("ledge: unexpected data size: %d vs %d - %s", len(data), size, string(slice))
	}
	return data, nil
}
