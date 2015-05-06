package ledge

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

var (
	rpcDecoderInstance = &rpcDecoder{}
)

type rpcDecoder struct{}

func (r *rpcDecoder) Decode(reader io.Reader) ([]byte, error) {
	sizeBytesLenBuf, err := r.read(reader, 1)
	if err != nil {
		return nil, err
	}
	sizeBytes, err := r.read(reader, int(sizeBytesLenBuf[0]))
	if err != nil {
		return nil, err
	}
	size, err := strconv.ParseInt(string(sizeBytes), 10, 64)
	if err != nil {
		return nil, err
	}
	sizeInt := int(size)
	if int64(sizeInt) != size {
		return nil, fmt.Errorf("ledge: could not cast %d to an int", size)
	}
	data, err := r.read(reader, sizeInt)
	if err != nil {
		return nil, err
	}
	// docker logs do not flush without a newline, this is the hack
	separatorBuf, err := r.read(reader, 1)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(separatorBuf, separatorBytes) {
		return nil, fmt.Errorf("ledge: expected separator (%v), %v, got (%v) %v", separatorBytes, string(separatorBytes), separatorBuf, string(separatorBuf))
	}
	return data, nil
}

func (r *rpcDecoder) read(reader io.Reader, size int) ([]byte, error) {
	buffer := make([]byte, size)
	readSoFar := 0
	for readSoFar < size {
		data := make([]byte, size-readSoFar)
		n, err := reader.Read(data)
		if err != nil {
			return nil, err
		}
		if n > 0 {
			m := copy(buffer[readSoFar:readSoFar+n], data[:n])
			if m != n {
				return nil, fmt.Errorf("ledge: tried to copy %d bytes, only copied %d bytes", n, m)
			}
			readSoFar += n
		}
	}
	return buffer, nil
}
