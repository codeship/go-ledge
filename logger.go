package ledge

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"time"
)

type logger struct {
	writer              io.Writer
	marshaller          Marshaller
	reflectTypeProvider *reflectTypeProvider
	options             LoggerOptions
	contexts            []Context
}

func newLogger(
	writer io.Writer,
	marshaller Marshaller,
	reflectTypeProvider *reflectTypeProvider,
	opts LoggerOptions,
	contexts []Context,
) *logger {
	return &logger{
		writer,
		marshaller,
		reflectTypeProvider,
		opts,
		contexts,
	}
}

func (l *logger) WithContext(context Context) Logger {
	if err := l.reflectTypeProvider.validateContextReflectType(reflect.TypeOf(context)); err != nil {
		panic(err.Error())
	}
	return newLogger(
		l.writer,
		l.marshaller,
		l.reflectTypeProvider,
		l.options,
		append(l.contexts, context),
	)
}

func (l *logger) Unstructured() UnstructuredLogger {
	return newUnstructuredLogger(l, make(map[string]interface{}))
}

func (l *logger) Debug(event Event) {
	l.print(LevelDebug, event)
}

func (l *logger) Error(event Event) {
	l.print(LevelError, event)
}

func (l *logger) Fatal(event Event) {
	l.print(LevelFatal, event)
}

func (l *logger) Info(event Event) {
	l.print(LevelInfo, event)
}

func (l *logger) Panic(event Event) {
	l.print(LevelPanic, event)
}

func (l *logger) Warn(event Event) {
	l.print(LevelWarn, event)
}

func (l *logger) DebugWriter(event Event) io.Writer {
	return l.printWriter(LevelDebug, event)
}

func (l *logger) ErrorWriter(event Event) io.Writer {
	return l.printWriter(LevelError, event)
}

func (l *logger) InfoWriter(event Event) io.Writer {
	return l.printWriter(LevelInfo, event)
}

func (l *logger) WarnWriter(event Event) io.Writer {
	return l.printWriter(LevelWarn, event)
}

func (l *logger) print(level Level, event Event) {
	if err := l.reflectTypeProvider.validateEventReflectType(reflect.TypeOf(event)); err != nil {
		panic(err.Error())
	}
	_, err := l.write(l.getEntry(l.getBaseEntry(level, event), nil))
	if err != nil {
		panic(err.Error())
	}
}

func (l *logger) printWriter(level Level, event Event) io.Writer {
	if err := l.reflectTypeProvider.validateEventReflectType(reflect.TypeOf(event)); err != nil {
		panic(err.Error())
	}
	return newEntryWriter(l.getBaseEntry(level, event), l)
}

func (l *logger) getBaseEntry(level Level, event Event) *Entry {
	return &Entry{
		Level:    level,
		Contexts: l.contexts,
		Event:    event,
	}
}

func (l *logger) getEntry(baseEntry *Entry, writerOutput []byte) *Entry {
	return &Entry{
		ID:           l.allocateID(),
		Time:         l.now(),
		Level:        baseEntry.Level,
		Contexts:     baseEntry.Contexts,
		Event:        baseEntry.Event,
		WriterOutput: writerOutput,
	}
}

func (l *logger) allocateID() string {
	if l.options.IDAllocator != nil {
		return l.options.IDAllocator.Allocate()
	}
	return uuidAllocatorInstance.Allocate()
}

func (l *logger) now() time.Time {
	if l.options.Timer != nil {
		return l.options.Timer.Now()
	}
	return systemTimerInstance.Now()
}

func (l *logger) write(entry *Entry) (int, error) {
	p, err := l.marshaller.Marshal(entry)
	if err != nil {
		return 0, err
	}
	// TODO(pedge): does this work?
	if entry.Level == LevelPanic {
		panic(string(p))
	}
	if l.include(entry) {
		if l.options.Encoder != nil {
			return l.options.Encoder.Encode(l.writer, p)
		}
		q, err := l.addNewline(p)
		if err != nil {
			return 0, err
		}
		return l.writer.Write(q)
	}
	// TODO(pedge): does this work?
	if entry.Level == LevelFatal {
		os.Exit(1)
	}
	return 0, nil
}

func (l *logger) include(entry *Entry) bool {
	return includeEntry(l.options.Filters, entry)
}

func (l *logger) addNewline(p []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	if _, err := buffer.Write(p); err != nil {
		return nil, err
	}
	if _, err := buffer.WriteRune('\n'); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
