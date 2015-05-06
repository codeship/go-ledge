package ledge

import (
	"io/ioutil"
	"testing"
	"time"
)

var (
	//testWriter        = io.Writer(os.Stdout)
	testWriter        = ioutil.Discard
	testSpecification = &Specification{
		ContextTypes: []Context{
			TestRequestID(""),
			TestInteger(0),
			TestContextBar{},
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
	One string
	Two int
}

func TestLoggerPrintToStdout(t *testing.T) {
	for _, marshaller := range []Marshaller{
		NewTextMarshaller(
			TextMarshallerOptions{},
		),
		ShortJSONMarshaller,
		JSONMarshaller,
	} {
		logger, err := NewLogger(
			testWriter,
			marshaller,
			testSpecification,
			LoggerOptions{
				WriteNewline: true,
			},
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

	if err := fakeLogger.CheckEntriesEqual(
		[]*Entry{
			&Entry{
				ID:    "0",
				Time:  time.Unix(0, 0),
				Level: LevelInfo,
				Event: TestEventFoo{"one", 2},
			},
			&Entry{
				ID:    "1",
				Time:  time.Unix(100, 0),
				Level: LevelInfo,
				Event: &TestEventFooPtr{"one", 2},
			},
			&Entry{
				ID:    "2",
				Time:  time.Unix(200, 0),
				Level: LevelInfo,
				Contexts: []Context{
					TestRequestID("bar"),
				},
				Event: TestEventFoo{"one", 2},
			},
			&Entry{
				ID:    "3",
				Time:  time.Unix(300, 0),
				Level: LevelInfo,
				Contexts: []Context{
					TestInteger(10),
				},
				Event: &TestEventFooPtr{"one", 2},
			},
			&Entry{
				ID:    "4",
				Time:  time.Unix(400, 0),
				Level: LevelInfo,
				Contexts: []Context{
					TestContextBar{"one", 2},
				},
				Event: Info("hello"),
			},
			&Entry{
				ID:    "5",
				Time:  time.Unix(500, 0),
				Level: LevelInfo,
				Contexts: []Context{
					TestContextBar{"one", 2},
				},
				Event: Info("{key:value} hello"),
			},
		},
		true,
		true,
	); err != nil {
		t.Error(err)
	}
}
