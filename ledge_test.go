package ledge

import (
	"bytes"
	"io/ioutil"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
)

var (
	//testWriter = io.Writer(os.Stdout)
	testWriter        = ioutil.Discard
	testSpecification = &Specification{
		ContextTypes: []Context{
			TestRequestID(""),
			TestInteger(0),
			TestContextBar{},
			Level_NONE,
		},
		EventTypes: []Event{
			TestEventFoo{},
			&TestEventFooPtr{},
		},
	}
)

type TestRequestID string

type TestInteger int

type TestContextBar struct {
	One string
	Two int
}

type TestEventFoo struct {
	One string
	Two int
}

type TestEventFooPtr struct {
	One string `protobuf:"bytes,1,opt,name=one" json:"one,omitempty"`
	Two int32  `protobuf:"varint,2,opt,name=two" json:"two,omitempty"`
}

func (m *TestEventFooPtr) Reset()         { *m = TestEventFooPtr{} }
func (m *TestEventFooPtr) String() string { return proto.CompactTextString(m) }
func (*TestEventFooPtr) ProtoMessage()    {}

func TestLoggerPrintToStdout(t *testing.T) {
	for _, marshaller := range []Marshaller{
		NewTextMarshaller(
			TextMarshallerOptions{},
		),
		JSONMarshaller,
	} {
		logger, err := NewLogger(
			testWriter,
			marshaller,
			testSpecification,
			LoggerOptions{},
		)
		if err != nil {
			t.Fatal(err)
		}
		logger.Info(TestEventFoo{"one", 2})
		logger.Info(&TestEventFooPtr{"one", 2})
		logger.WithContext(TestRequestID("bar")).Info(TestEventFoo{"one", 2})
		logger.WithContext(TestInteger(10)).Info(&TestEventFooPtr{"one", 2})
		logger.WithContext(TestContextBar{"one", 2}).Unstructured().Info("hello")
		logger.WithContext(TestContextBar{"one", 2}).Unstructured().WithField("key", "value").Info("hello")
		logger.Unstructured().WithField("key", "value").Info("")
		logger.Unstructured().Info("")
		logger.WithContext(Level_PANIC).Unstructured().Info("panic context")
	}
}

func TestFakeLogger(t *testing.T) {
	fakeLogger, err := NewFakeLogger(testSpecification)
	if err != nil {
		t.Fatal(err)
	}
	fakeLogger.Info(TestEventFoo{"one", 2})
	fakeLogger.AddTimeSec(100)
	fakeLogger.Info(&TestEventFooPtr{"one", 2})
	fakeLogger.AddTimeSec(100)
	fakeLogger.WithContext(TestRequestID("bar")).Info(TestEventFoo{"one", 2})
	fakeLogger.AddTimeSec(100)
	fakeLogger.WithContext(TestInteger(10)).Info(&TestEventFooPtr{"one", 2})
	fakeLogger.AddTimeSec(100)
	fakeLogger.WithContext(TestContextBar{"one", 2}).Unstructured().Info("hello")
	fakeLogger.AddTimeSec(100)
	fakeLogger.WithContext(TestContextBar{"one", 2}).Unstructured().WithField("key", "value").Info("hello")
	fakeLogger.AddTimeSec(100)
	fakeLogger.Unstructured().WithField("key", "value").Info("")
	fakeLogger.AddTimeSec(100)
	fakeLogger.Unstructured().Info("")
	fakeLogger.AddTimeSec(100)
	fakeLogger.WithContext(Level_PANIC).Unstructured().Info("panic context")

	if err := fakeLogger.CheckEntriesEqual(
		[]*Entry{
			&Entry{
				ID:    "0",
				Time:  time.Unix(0, 0),
				Level: Level_INFO,
				Event: TestEventFoo{"one", 2},
			},
			&Entry{
				ID:    "1",
				Time:  time.Unix(100, 0),
				Level: Level_INFO,
				Event: &TestEventFooPtr{"one", 2},
			},
			&Entry{
				ID:    "2",
				Time:  time.Unix(200, 0),
				Level: Level_INFO,
				Contexts: []Context{
					TestRequestID("bar"),
				},
				Event: TestEventFoo{"one", 2},
			},
			&Entry{
				ID:    "3",
				Time:  time.Unix(300, 0),
				Level: Level_INFO,
				Contexts: []Context{
					TestInteger(10),
				},
				Event: &TestEventFooPtr{"one", 2},
			},
			&Entry{
				ID:    "4",
				Time:  time.Unix(400, 0),
				Level: Level_INFO,
				Contexts: []Context{
					TestContextBar{"one", 2},
				},
				Event: &UnstructuredEvent{"hello"},
			},
			&Entry{
				ID:    "5",
				Time:  time.Unix(500, 0),
				Level: Level_INFO,
				Contexts: []Context{
					TestContextBar{"one", 2},
				},
				Event: &UnstructuredEvent{"{key:value} hello"},
			},
			&Entry{
				ID:    "6",
				Time:  time.Unix(600, 0),
				Level: Level_INFO,
				Event: &UnstructuredEvent{"{key:value}"},
			},
			&Entry{
				ID:    "7",
				Time:  time.Unix(700, 0),
				Level: Level_INFO,
				Event: &UnstructuredEvent{""},
			},
			&Entry{
				ID:    "8",
				Time:  time.Unix(800, 0),
				Level: Level_INFO,
				Contexts: []Context{
					Level_PANIC,
				},
				Event: &UnstructuredEvent{"panic context"},
			},
		},
		true,
		true,
	); err != nil {
		t.Error(err)
	}
}

func TestRoundTrip(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	timer := newFakeTimer(0)
	logger, err := NewLogger(
		buffer,
		ProtoMarshaller,
		testSpecification,
		LoggerOptions{
			IDAllocator: newFakeIDAllocator(),
			Timer:       timer,
			Encoder:     RPCEncoder,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	logger.Info(TestEventFoo{"one", 2})
	timer.AddTimeSec(100)
	logger.Info(&TestEventFooPtr{"one", 2})
	timer.AddTimeSec(100)
	logger.WithContext(TestRequestID("bar")).Info(TestEventFoo{"one", 2})
	timer.AddTimeSec(100)
	logger.WithContext(TestInteger(10)).Info(&TestEventFooPtr{"one", 2})
	timer.AddTimeSec(100)
	logger.WithContext(TestContextBar{"one", 2}).Unstructured().Info("hello")
	timer.AddTimeSec(100)
	logger.WithContext(TestContextBar{"one", 2}).Unstructured().WithField("key", "value").Info("hello")
	timer.AddTimeSec(100)
	logger.Unstructured().WithField("key", "value").Info("")
	timer.AddTimeSec(100)
	logger.Unstructured().Info("")
	timer.AddTimeSec(100)
	byteString := string([]byte{128})
	logger.Unstructured().Info(byteString)

	unmarshaller, err := NewProtoUnmarshaller(testSpecification)
	if err != nil {
		t.Fatal(err)
	}
	entryReader, err := NewEntryReader(
		buffer,
		unmarshaller,
		RPCDecoder,
		EntryReaderOptions{},
	)
	if err != nil {
		t.Fatal(err)
	}
	entries, err := NewBlockingEntryReader(entryReader).Entries()
	if err != nil {
		t.Fatal(err)
	}

	if err := checkEntriesEqual(
		entries,
		[]*Entry{
			&Entry{
				ID:    "0",
				Time:  time.Unix(0, 0),
				Level: Level_INFO,
				Event: TestEventFoo{"one", 2},
			},
			&Entry{
				ID:    "1",
				Time:  time.Unix(100, 0),
				Level: Level_INFO,
				Event: &TestEventFooPtr{"one", 2},
			},
			&Entry{
				ID:    "2",
				Time:  time.Unix(200, 0),
				Level: Level_INFO,
				Contexts: []Context{
					TestRequestID("bar"),
				},
				Event: TestEventFoo{"one", 2},
			},
			&Entry{
				ID:    "3",
				Time:  time.Unix(300, 0),
				Level: Level_INFO,
				Contexts: []Context{
					TestInteger(10),
				},
				Event: &TestEventFooPtr{"one", 2},
			},
			&Entry{
				ID:    "4",
				Time:  time.Unix(400, 0),
				Level: Level_INFO,
				Contexts: []Context{
					TestContextBar{"one", 2},
				},
				Event: &UnstructuredEvent{"hello"},
			},
			&Entry{
				ID:    "5",
				Time:  time.Unix(500, 0),
				Level: Level_INFO,
				Contexts: []Context{
					TestContextBar{"one", 2},
				},
				Event: &UnstructuredEvent{"{key:value} hello"},
			},
			&Entry{
				ID:    "6",
				Time:  time.Unix(600, 0),
				Level: Level_INFO,
				Event: &UnstructuredEvent{"{key:value}"},
			},
			&Entry{
				ID:    "7",
				Time:  time.Unix(700, 0),
				Level: Level_INFO,
				Event: &UnstructuredEvent{""},
			},
			&Entry{
				ID:    "8",
				Time:  time.Unix(800, 0),
				Level: Level_INFO,
				Event: &UnstructuredEvent{byteString},
			},
		},
		true,
		true,
	); err != nil {
		t.Error(err)
	}
}
