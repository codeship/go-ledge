package ledge

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
)

const (
	// TODO(pedge): convert to unix nanos?
	timeFormat = "2006-01-02 15:04:05.999999999 -0700 MST"
)

var (
	defaultJSONKeys = &jsonKeys{
		id:           "id",
		time:         "time",
		level:        "log_level",
		eventType:    "event_type",
		writerOutput: "writer_output",
	}
	protoMarshallerInstance = &protoMarshaller{}
)

type textMarshaller struct {
	options TextMarshallerOptions
}

func newTextMarshaller(
	options TextMarshallerOptions,
) *textMarshaller {
	return &textMarshaller{
		options,
	}
}

func (t *textMarshaller) Marshal(entry *Entry) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	// can't use [ because zsh complains
	if _, err := buffer.WriteString("{"); err != nil {
		return nil, err
	}
	if !t.options.NoID {
		if _, err := buffer.WriteString(entry.ID); err != nil {
			return nil, err
		}
		if _, err := buffer.WriteString(" "); err != nil {
			return nil, err
		}
	}
	if !t.options.NoTime {
		if _, err := buffer.WriteString(entry.Time.Format("15:04:05.000000000")); err != nil {
			return nil, err
		}
		if _, err := buffer.WriteString(" "); err != nil {
			return nil, err
		}
	}
	if !t.options.NoLevel {
		if _, err := buffer.WriteString(strings.ToLower(entry.Level.String())); err != nil {
			return nil, err
		}
		if _, err := buffer.WriteString(" "); err != nil {
			return nil, err
		}
	}
	for _, context := range entry.Contexts {
		contextString, err := t.objectString(context)
		if err != nil {
			return nil, err
		}
		if _, err := buffer.WriteString(contextString); err != nil {
			return nil, err
		}
		if _, err := buffer.WriteString(" "); err != nil {
			return nil, err
		}
	}
	eventString, err := t.objectString(entry.Event)
	if err != nil {
		return nil, err
	}
	if _, err := buffer.WriteString(eventString); err != nil {
		return nil, err
	}
	if _, err := buffer.WriteString("}"); err != nil {
		return nil, err
	}
	if entry.WriterOutput != nil && len(entry.WriterOutput) != 0 {
		if _, err := buffer.WriteString(": "); err != nil {
			return nil, err
		}
		if _, err := buffer.Write(entry.WriterOutput); err != nil {
			return nil, err
		}
	}
	return []byte(strings.TrimSpace(buffer.String())), nil
}

func (t *textMarshaller) objectString(object interface{}) (string, error) {
	if object == nil {
		return "nil", nil
	}
	reflectType := reflect.TypeOf(object)
	typeString, err := shortReflectKey(reflectType)
	if err != nil {
		return "", err
	}
	if stringer, ok := object.(fmt.Stringer); ok {
		return fmt.Sprintf("%s=%s", typeString, stringer.String()), nil
	}
	objectString := fmt.Sprintf("%+v", object)
	if len(objectString) > 0 && objectString[0:1] == "&" {
		objectString = objectString[1:]
	}
	if len(objectString) > 1 && objectString[0:2] == "{ " {
		objectString = fmt.Sprintf("{%s", objectString[2:])
	}
	return fmt.Sprintf("%s=%s", typeString, objectString), nil
}

type jsonKeys struct {
	id           string
	time         string
	level        string
	eventType    string
	writerOutput string
}

type jsonMarshaller struct {
	jsonKeys *jsonKeys
}

func newJSONMarshaller(
	jsonKeys *jsonKeys,
) *jsonMarshaller {
	return &jsonMarshaller{
		jsonKeys,
	}
}

func (j *jsonMarshaller) Marshal(entry *Entry) ([]byte, error) {
	m := make(map[string]interface{})
	m[j.jsonKeys.id] = entry.ID
	m[j.jsonKeys.time] = entry.Time.Format(timeFormat)
	m[j.jsonKeys.level] = strings.ToLower(entry.Level.String())
	for _, context := range entry.Contexts {
		contextKey, err := shortReflectKey(reflect.TypeOf(context))
		if err != nil {
			return nil, err
		}
		m[contextKey] = context
	}
	eventKey, err := shortReflectKey(reflect.TypeOf(entry.Event))
	if err != nil {
		return nil, err
	}
	m[j.jsonKeys.eventType] = eventKey
	m[eventKey] = entry.Event
	m[j.jsonKeys.writerOutput] = string(entry.WriterOutput)
	return json.Marshal(m)
}

func shortReflectKey(reflectType reflect.Type) (string, error) {
	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	name := reflectType.Name()
	if name == "" {
		return "", fmt.Errorf("ledge: no name for type %v", reflectType)
	}
	return name, nil
}

type protoMarshaller struct{}

func (p *protoMarshaller) Marshal(entry *Entry) ([]byte, error) {
	protoEntry := &ProtoEntry{
		Id:           entry.ID,
		TimeUnixNsec: entry.Time.UnixNano(),
		Level:        entry.Level,
		ContextTypeNameToContext: make(map[string][]byte),
		WriterOutput:             entry.WriterOutput,
	}
	eventTypeName, err := getReflectTypeName(reflect.TypeOf(entry.Event))
	if err != nil {
		return nil, err
	}
	eventBytes, err := marshalBinary(entry.Event)
	if err != nil {
		return nil, err
	}
	protoEntry.EventTypeName = eventTypeName
	protoEntry.Event = eventBytes
	for _, context := range entry.Contexts {
		contextTypeName, err := getReflectTypeName(reflect.TypeOf(context))
		if err != nil {
			return nil, err
		}
		contextBytes, err := marshalBinary(context)
		if err != nil {
			return nil, err
		}
		protoEntry.ContextTypeNameToContext[contextTypeName] = contextBytes
	}
	b, err := proto.Marshal(protoEntry)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(nil)
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	if _, err := encoder.Write(b); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	bufferBytes := buffer.Bytes()
	return bufferBytes, nil
}
