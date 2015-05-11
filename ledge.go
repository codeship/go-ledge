package ledge

import (
	"io"
	"time"
)

var (
	// DebugFilter is a Filter that only includes Entries with a Level of at least LevelDebug.
	DebugFilter Filter = newLevelFilter(LevelDebug)
	// InfoFilter is a Filter that only includes Entries with a Level of at least LevelInfo.
	InfoFilter Filter = newLevelFilter(LevelInfo)
	// WarnFilter is a Filter that only includes Entries with a Level of at least LevelWarn.
	WarnFilter Filter = newLevelFilter(LevelWarn)
	// ErrorFilter is a Filter that only includes Entries with a Level of at least LevelError.
	ErrorFilter Filter = newLevelFilter(LevelError)
	// FatalFilter is a Filter that only includes Entries with a Level of at least LevelFatal.
	FatalFilter Filter = newLevelFilter(LevelFatal)
	// PanicFilter is a Filter that only includes Entries with a Level of at least LevelPanic.
	PanicFilter Filter = newLevelFilter(LevelPanic)

	// ShortJSONMarshaller is a Marshaller that marshales Entries in JSON, but with shorthand
	// notation for Context and Entry types. It should not be used for logging intended for RPC use.
	ShortJSONMarshaller Marshaller = newShortJSONMarshaller(defaultJSONKeys)
	// JSONMarshaller is a Marshaller for JSON. It is intended for RPC use.
	JSONMarshaller Marshaller = newJSONMarshaller(defaultJSONKeys)
	// RPCEncoder is an Encoder that wraps data in a simple RPC format.
	RPCEncoder Encoder = rpcEncoderInstance
	// RPCDecoder is a Decoder that decodes data encoded with RPCEncoder.
	RPCDecoder Decoder = rpcDecoderInstance

	// DefaultEventTypes are the Event types included with every Logger, EntryReader,
	// and BlockingEntryReader by default. These are used for the UnstructuredLogger.
	DefaultEventTypes = []Event{
		Debug(""),
		Error(""),
		Fatal(""),
		Info(""),
		Panic(""),
		Warn(""),
	}
)

// A Context is attached to a Logger and included as part of every Entry a Logger outputs.
type Context interface{}

// An Event is outputted by a Logger.
type Event interface{}

// Debug is the Event for Debug statements with an UnstructuredLogger.
type Debug string

// Error is the Event for Error statements with an UnstructuredLogger.
type Error string

// Fatal is the Event for Fatal statements with an UnstructuredLogger.
type Fatal string

// Info is the Event for Info statements with an UnstructuredLogger.
type Info string

// Panic is the Event for Panic statements with an UnstructuredLogger.
type Panic string

// Warn is the Event for Warn statements with an UnstructuredLogger.
type Warn string

// Fields are attached to an UnstructuredLogger and included as part of every statement outputted.
type Fields map[string]interface{}

// An UnstructuredLogger allows logging without the use of typed Events. This is meant to be used
// for quick additional logger, for adoption, and as a replacement for Golang's standard logger.
// In general, if using the idioms of this library, use of UnstructuredLogger is discouraged.
type UnstructuredLogger interface {
	// WithField returns a new UnstructuredLogger with the given field value.
	WithField(key string, value interface{}) UnstructuredLogger
	// WithField returns a new UnstructuredLogger with the given Fields.
	WithFields(fields Fields) UnstructuredLogger

	// Debug prints a Debug event, analogous to fmt.Sprint.
	Debug(args ...interface{})
	// Debugf prints a Debug event, analogous to fmt.Sprintf.
	Debugf(format string, args ...interface{})
	// Debugln prints a Debug event, analogous to fmt.Sprintln.
	Debugln(args ...interface{})
	// Error prints an Error event, analogous to fmt.Sprint.
	Error(args ...interface{})
	// Errorf prints an Error event, analogous to fmt.Sprintf.
	Errorf(format string, args ...interface{})
	// Errorln prints an Error event, analogous to fmt.Sprintln.
	Errorln(args ...interface{})
	// Fatal prints a Fatal event, analogous to fmt.Sprint. It then exits with os.Exit(1).
	Fatal(args ...interface{})
	// Fatalf prints a Fatal event, analogous to fmt.Sprintf. It then exits with os.Exit(1).
	Fatalf(format string, args ...interface{})
	// Fatalln prints a Fatal event, analogous to fmt.Sprintln. It then exits with os.Exit(1).
	Fatalln(args ...interface{})
	// Info prints an Info event, analogous to fmt.Sprint.
	Info(args ...interface{})
	// Infof prints an Info event, analogous to fmt.Sprintf.
	Infof(format string, args ...interface{})
	// Infoln prints an Info event, analogous to fmt.Sprintln.
	Infoln(args ...interface{})
	// Panic with a Panic event, analogous to fmt.Sprint.
	Panic(args ...interface{})
	// Panic with a Panic event, analogous to fmt.Sprintf.
	Panicf(format string, args ...interface{})
	// Panic a Panic event, analogous to fmt.Sprintln.
	Panicln(args ...interface{})
	// Print is an alias for Info.
	Print(args ...interface{})
	// Printf is an alias for Infof.
	Printf(format string, args ...interface{})
	// Println is an alias for Infoln.
	Println(args ...interface{})
	// Warn prints a Warn event, analogous to fmt.Sprint.
	Warn(args ...interface{})
	// Warnf prints a Warn event, analogous to fmt.Sprintf.
	Warnf(format string, args ...interface{})
	// Warnln prints a Warn event, analogous to fmt.Sprintln.
	Warnln(args ...interface{})

	// DebugWriter returns a new io.Writer that will log output to the writer
	// inside the WriterOutput field of an Entry at the Debug Level.
	DebugWriter() io.Writer
	// ErrorWriter returns a new io.Writer that will log output to the writer
	// inside the WriterOutput field of an Entry at the Error Level.
	ErrorWriter() io.Writer
	// InfoWriter returns a new io.Writer that will log output to the writer
	// inside the WriterOutput field of an Entry at the Info Level.
	InfoWriter() io.Writer
	// WarnWriter returns a new io.Writer that will log output to the writer
	// inside the WriterOutput field of an Entry at the Warn Level.
	WarnWriter() io.Writer
}

// Logger is the main logging interface. A Logger logs Events with given Contexts as Entry objects.
type Logger interface {
	// WithContext returns a new Logger with the given Context attached. If the Context
	// was not registered in the Specification on Logger creation, this method will panic.
	WithContext(context Context) Logger
	// Unstructured returns the associated UnstructuredLogger. The methods on UnstructuredLogger
	// are not directly included on Logger to discourage use of these methods.
	Unstructured() UnstructuredLogger

	// Debug prints an event at the Debug Level.
	Debug(event Event)
	// Error prints an event at the Error Level.
	Error(event Event)
	// Fatal prints an event at the Fatal Level.
	Fatal(event Event)
	// Info prints an event at the Info Level.
	Info(event Event)
	// Panic prints an event at the Panic Level.
	Panic(event Event)
	// Warn prints an event at the Warn Level.
	Warn(event Event)
	// DebugWriter returns a new io.Writer that will log output to the writer
	// inside the WriterOutput field of an Entry, using the associated Event,
	// at the Debug Level.
	DebugWriter(event Event) io.Writer
	// ErrorWriter returns a new io.Writer that will log output to the writer
	// inside the WriterOutput field of an Entry, using the associated Event,
	// at the Error Level.
	ErrorWriter(event Event) io.Writer
	// InfoWriter returns a new io.Writer that will log output to the writer
	// inside the WriterOutput field of an Entry, using the associated Event,
	// at the Info Level.
	InfoWriter(event Event) io.Writer
	// WarnWriter returns a new io.Writer that will log output to the writer
	// inside the WriterOutput field of an Entry, using the associated Event,
	// at the Warn Level.
	WarnWriter(event Event) io.Writer
}

// Entry is the type that is marshalled and unmarshalled into and from log messages.
// Every log message is an Entry, including UnstructuredLogger messages.
type Entry struct {
	// ID is a unique ID that is assocated with this Entry, usually a UUID except in testing.
	// This must be globally unique across multiple instances of a Logger.
	ID string
	// Time is the time that this Entry was logged.
	Time time.Time
	// Level is the Level of this entry.
	Level Level
	// Contexts is the contexts that were associated with the Logger when this entry was logged.
	Contexts []Context
	// Event is the event that was logged.
	Event Event
	// WriterOutput is the associated writer output, if this Entry was used for a Writer function.
	// If this Entry was created from a non-Writer function, this will be nil.
	WriterOutput []byte
}

// IDAllocator allocated unique IDs for Entry structs.
type IDAllocator interface {
	// Allocate allocates a new unique ID.
	Allocate() string
}

// Timer gets the current time.
type Timer interface {
	// Now returns the current time.
	Now() time.Time
}

// Filter allows filtering of Entry objects when reading or writing.
type Filter interface {
	// Include returns true if the Entry should be included.
	Include(entry *Entry) bool
}

// NewRequireContextFilter returns a filter that only selects Entry objects with the given Context.
func NewRequireContextFilter(context Context) Filter {
	return newRequireContextFilter(
		context,
	)
}

// Marshaller marshals Entry objects into byte slices.
type Marshaller interface {
	// Marshal marshals Entry objects into byte slices.
	Marshal(entry *Entry) ([]byte, error)
}

// TextMarshallerOptions provides options for creating TextMarshallers.
type TextMarshallerOptions struct {
	// NoID will suppress the printing of Entry IDs.
	NoID bool
	// NoTime will suppress the printing of Entry times.
	NoTime bool
	// NoLevel will suppress the printing of Entry Levels.
	NoLevel bool
}

// NewTextMarshaller returns a Marshaller that marshals output in a human-readable manner.
// This should never be used if an EntryReader or BlockingEntryReader is to be used with the Entry objects.
func NewTextMarshaller(options TextMarshallerOptions) Marshaller {
	return newTextMarshaller(
		options,
	)
}

// Encoder encodes marshalled byte slices to a writer, optionally adding output.
type Encoder interface {
	// Encode encodes marshalled byte slices to a writer, optionally adding output.
	Encode(writer io.Writer, p []byte) (int, error)
}

// Specification specifies the Context and Event types that will be used with a Logger, EntryReader,
// or BlockingEntryReader. A type is specified using the zero value. For example, given:
//
//		type RequestId string
//		type FooEvent struct {
//			One string
//		}
//
// And assuming a FooEvent is used as a pointer, the Specification should be:
//
//		var (
//			AppSpecifiation = &Specification{
//				ContextTypes: []Context{
//					RequestId(""),
//				},
//				EventTypes: []Event{
//					&FooEvent{},
//				},
//			}
//		)
type Specification struct {
	ContextTypes []Context
	EventTypes   []Event
}

// MergeSpecifications merges multiple Specifications into a single specification.
func MergeSpecifications(specifications ...*Specification) *Specification {
	return mergeSpecifications(specifications)
}

// LoggerOptions specifies the options to be used when creating a Logger.
type LoggerOptions struct {
	// IDAllocator specifies an alternate IDAllocator to use.
	// If not specified, a UUID allocator will be used.
	IDAllocator IDAllocator
	// Time specifies an alternate Timer to use.
	// If not specification, a system Timer will be used.
	Timer Timer
	// Filters specifies the Filters to use.
	Filters []Filter
	// Encoder specifies an Encoder to use.
	// If not specified, no encoder will be used and marshalled Entry objects
	// will be directed printed to the Logger's io.Writer.
	Encoder Encoder
	// WriteNewline specifies whether to write a '\n' at the end of a marshalled Entry Object.
	WriteNewline bool
}

// NewLogger creates a new Logger.
func NewLogger(writer io.Writer, marshaller Marshaller, specification *Specification, options LoggerOptions) (Logger, error) {
	reflectTypeProvider, err := newReflectTypeProvider(specification)
	if err != nil {
		return nil, err
	}
	return newLogger(
		writer,
		marshaller,
		reflectTypeProvider,
		options,
		make([]Context, 0),
	), nil
}

// Unmarshaller unmarshals a byte slice into an Entry.
type Unmarshaller interface {
	// Unmarshal unmarshals a byte slice into an Entry.
	Unmarshal(p []byte) (*Entry, error)
}

// NewJSONUnmarshaller returns a new Unmarshaller that unmarshals Entry objects
// marshalled with JSONMarshaller.
func NewJSONUnmarshaller(specification *Specification) (Unmarshaller, error) {
	return newJSONUnmarshaller(
		specification,
		defaultJSONKeys,
	)
}

// Decoder decodes an input stream into separate byte slices that represent marshalled Entry objects.
type Decoder interface {
	// Decode gets the next marshalled Entry object from the input stream.
	Decode(reader io.Reader) ([]byte, error)
}

// EntryResponse is a response from an EntryReader.
type EntryResponse struct {
	// Entry is the Entry read.
	Entry *Entry
	// Error will be set if there was an error reading.
	Error error
}

// EntryReader reads Entry objects from an input stream.
type EntryReader interface {
	// Channel returns a read channel of EntryResponse objects.
	Channel() <-chan *EntryResponse
	// Cancel cancels reading and will close the channel.
	Cancel() error
}

// EntryReaderOptions specifies the options to be used when creating an EntryReader.
type EntryReaderOptions struct {
	// Filters specifies the Filters to use.
	Filters []Filter
}

// NewEntryReader returns a new EntryReader.
func NewEntryReader(reader io.Reader, unmarshaller Unmarshaller, decoder Decoder, options EntryReaderOptions) (EntryReader, error) {
	return newEntryReader(
		reader,
		unmarshaller,
		decoder,
		options,
	)
}

// BlockingEntryReader reads Entry objects in a blocking manner until the input stream is finished.
type BlockingEntryReader interface {
	// Entries returns all Entry objects in the order they were read.
	Entries() ([]*Entry, error)
}

// NewBlockingEntryReader returns a new BlockingEntryReader.
func NewBlockingEntryReader(entryReader EntryReader) BlockingEntryReader {
	return newBlockingEntryReader(
		entryReader,
	)
}

// FakeLogger is a Logger and BlockingEntryReader that can be used to test log output.
// It uses a fake IDAllocator that allocates IDs as integers, starting from 0, and a fake
// Timer that starts at unix time 0. Entry objects are logged to an internal buffer,
// which can be read using the BlockingEntryReader functionality.
type FakeLogger interface {
	Logger
	BlockingEntryReader

	// AddTimeSec adds the specified number of seconds to the fake Timer.
	AddTimeSec(int64)

	// CheckEntriesEqual returns error if the expected Entry objects do not match the logged Entry objects.
	CheckEntriesEqual(expected []*Entry, checkID bool, checkTime bool) error
}

// NewFakeLogger returns a new FakeLogger.
func NewFakeLogger(specification *Specification) (FakeLogger, error) {
	return newFakeLogger(
		specification,
	)
}
