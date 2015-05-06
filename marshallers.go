package ledge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
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
	shortJSONReflectTypeKeyProviderInstance = &shortJSONReflectTypeKeyProvider{}
	longJSONReflectTypeKeyProviderInstance  = &longJSONReflectTypeKeyProvider{}
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
		if _, err := buffer.WriteString(entry.Level.String()); err != nil {
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
	typeString, err := shortJSONReflectTypeKeyProviderInstance.Key(reflectType)
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

type jsonReflectTypeKeyProvider interface {
	Key(reflectType reflect.Type) (string, error)
}

type shortJSONReflectTypeKeyProvider struct{}

func (s *shortJSONReflectTypeKeyProvider) Key(reflectType reflect.Type) (string, error) {
	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	name := reflectType.Name()
	if name == "" {
		return "", fmt.Errorf("ledge: no name for type %v", reflectType)
	}
	return name, nil

}

type longJSONReflectTypeKeyProvider struct{}

func (l *longJSONReflectTypeKeyProvider) Key(reflectType reflect.Type) (string, error) {
	return getFullyQualifiedName(reflectType)
}

type jsonKeys struct {
	id           string
	time         string
	level        string
	eventType    string
	writerOutput string
}

type baseJSONMarshaller struct {
	jsonReflectTypeKeyProvider jsonReflectTypeKeyProvider
	jsonKeys                   *jsonKeys
}

func newShortJSONMarshaller(
	jsonKeys *jsonKeys,
) *baseJSONMarshaller {
	return newBaseJSONMarshaller(
		shortJSONReflectTypeKeyProviderInstance,
		jsonKeys,
	)
}

func newJSONMarshaller(
	jsonKeys *jsonKeys,
) *baseJSONMarshaller {
	return newBaseJSONMarshaller(
		longJSONReflectTypeKeyProviderInstance,
		jsonKeys,
	)
}

func newBaseJSONMarshaller(
	jsonReflectTypeKeyProvider jsonReflectTypeKeyProvider,
	jsonKeys *jsonKeys,
) *baseJSONMarshaller {
	if jsonKeys == nil {
		jsonKeys = defaultJSONKeys
	}
	return &baseJSONMarshaller{
		jsonReflectTypeKeyProvider,
		jsonKeys,
	}
}

func (b *baseJSONMarshaller) Marshal(entry *Entry) ([]byte, error) {
	m := make(map[string]interface{})
	m[b.jsonKeys.id] = entry.ID
	m[b.jsonKeys.time] = entry.Time.Format(timeFormat)
	m[b.jsonKeys.level] = entry.Level.String()
	for _, context := range entry.Contexts {
		contextKey, err := b.jsonReflectTypeKeyProvider.Key(reflect.TypeOf(context))
		if err != nil {
			return nil, err
		}
		m[contextKey] = context
	}
	eventKey, err := b.jsonReflectTypeKeyProvider.Key(reflect.TypeOf(entry.Event))
	if err != nil {
		return nil, err
	}
	m[b.jsonKeys.eventType] = eventKey
	m[eventKey] = entry.Event
	return json.Marshal(m)
}
