package ledge

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type fakeLogger struct {
	Logger
	BlockingEntryReader
	*fakeTimer
}

func newFakeLogger(
	specification *Specification,
) (*fakeLogger, error) {
	buffer := newLockedBuffer()
	unmarshaller, err := NewProtoUnmarshaller(specification)
	if err != nil {
		return nil, err
	}
	entryReader, err := NewEntryReader(
		buffer,
		unmarshaller,
		RPCDecoder,
		EntryReaderOptions{},
		100,
	)
	if err != nil {
		return nil, err
	}
	fakeIDAllocator := newFakeIDAllocator()
	fakeTimer := newFakeTimer(0)
	logger, err := NewLogger(
		buffer,
		ProtoMarshaller,
		specification,
		LoggerOptions{
			IDAllocator: fakeIDAllocator,
			Timer:       fakeTimer,
			Encoder:     RPCEncoder,
		},
	)
	if err != nil {
		return nil, err
	}
	return &fakeLogger{
		logger,
		NewBlockingEntryReader(entryReader),
		fakeTimer,
	}, nil
}

func (f *fakeLogger) CheckEntriesEqual(
	expected []*Entry,
	checkID bool,
	checkTime bool,
) error {
	entries, err := f.Entries()
	if err != nil {
		return err
	}
	return checkEntriesEqual(entries, expected, checkID, checkTime)
}

type fakeIDAllocator struct {
	value int64
}

func newFakeIDAllocator() *fakeIDAllocator {
	return &fakeIDAllocator{
		-1,
	}
}

func (ti *fakeIDAllocator) Allocate() string {
	return fmt.Sprintf("%d", atomic.AddInt64(&ti.value, 1))
}

type fakeTimer struct {
	now int64
}

func newFakeTimer(
	initialTimeUnixSec int64,
) *fakeTimer {
	return &fakeTimer{
		initialTimeUnixSec,
	}
}

func (tt *fakeTimer) AddTimeSec(delta int64) {
	tt.now += delta
}

func (tt *fakeTimer) Now() time.Time {
	return time.Unix(tt.now, 0)
}

type lockedBuffer struct {
	buffer *bytes.Buffer
	lock   *sync.Mutex
}

func newLockedBuffer() *lockedBuffer {
	return &lockedBuffer{
		bytes.NewBuffer(nil),
		&sync.Mutex{},
	}
}

func (l *lockedBuffer) Read(p []byte) (int, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.buffer.Read(p)
}

func (l *lockedBuffer) Write(p []byte) (int, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.buffer.Write(p)
}
