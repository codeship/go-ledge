package ledge

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

var (
	rpcEncoderInstance = &rpcEncoder{}

	separator                      = byte('\n')
	separatorBytes                 = []byte{separator}
	errShouldNotBeWritingZeroBytes = errors.New("ledge: should not be writing 0 bytes")
)

type rpcEncoder struct{}

func (r *rpcEncoder) Encode(writer io.Writer, data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}
	size := strconv.FormatInt(int64(len(data)), 10)
	sizeBytes := []byte(size)
	sizeBytesLen := len(sizeBytes)
	sizeBytesLenByte := byte(sizeBytesLen)
	if int(sizeBytesLenByte) != sizeBytesLen {
		return 0, fmt.Errorf("ledge: could not cast %d to a byte", sizeBytesLen)
	}
	if err := r.write(writer, []byte{sizeBytesLenByte}); err != nil {
		return 0, err
	}
	if err := r.write(writer, sizeBytes); err != nil {
		return 0, err
	}
	if err := r.write(writer, data); err != nil {
		return 0, err
	}
	// docker logs do not flush without a newline, this is the hack
	if err := r.writeSeparator(writer); err != nil {
		return 0, err
	}
	return len(data), nil
}

func (r *rpcEncoder) write(writer io.Writer, p []byte) error {
	if len(p) == 0 {
		return errShouldNotBeWritingZeroBytes
	}
	n, err := writer.Write(p)
	if err != nil {
		return err
	}
	if n != len(p) {
		return fmt.Errorf("ledge: tried to write %d bytes, wrote %d", len(p), n)
	}
	return nil
}

func (r *rpcEncoder) writeSeparator(writer io.Writer) error {
	return r.write(writer, separatorBytes)
}
