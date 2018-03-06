package ledge

import (
	"bufio"
	"io"
)

const (
	readerSize = 256 * 1024
)

type entryReader struct {
	reader       *bufio.Reader
	unmarshaller Unmarshaller
	decoder      Decoder
	options      EntryReaderOptions
	output       chan *EntryResponse
	cancel       chan bool
}

func newEntryReader(
	reader io.Reader,
	unmarshaller Unmarshaller,
	decoder Decoder,
	options EntryReaderOptions,
) (*entryReader, error) {
	obj := &entryReader{
		bufio.NewReaderSize(reader, readerSize),
		unmarshaller,
		decoder,
		options,
		make(chan *EntryResponse),
		make(chan bool),
	}
	go obj.read()
	return obj, nil
}

func (e *entryReader) Channel() <-chan *EntryResponse {
	return e.output
}

func (e *entryReader) Cancel() error {
	e.cancel <- true
	return nil
}

func (e *entryReader) read() {
	for {
		select {
		case <-e.cancel:
			close(e.output)
			return
		default:
			ok := e.keepReading()
			if !ok {
				close(e.output)
				return
			}
		}
	}
}

func (e *entryReader) keepReading() bool {
	data, err := e.decoder.Decode(e.reader)
	if err != nil {
		if err == io.EOF {
			return false
		}
		e.output <- &EntryResponse{Error: err}
		return true
	}
	entry, err := e.unmarshaller.Unmarshal(data)
	if err != nil {
		e.output <- &EntryResponse{Error: err}
		return true
	}
	if e.include(entry) {
		e.output <- &EntryResponse{Entry: entry}
	}
	return true
}

func (e *entryReader) include(entry *Entry) bool {
	return includeEntry(e.options.Filters, entry)
}
