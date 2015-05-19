package ledge

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
)

type jsonUnmarshaller struct {
	reflectTypeHandler *reflectTypeHandler
	jsonKeys           *jsonKeys
}

func newJSONUnmarshaller(
	specification *Specification,
	jsonKeys *jsonKeys,
) (*jsonUnmarshaller, error) {
	reflectTypeHandler, err := newReflectTypeHandler(specification)
	if err != nil {
		return nil, err
	}
	return &jsonUnmarshaller{
		reflectTypeHandler,
		jsonKeys,
	}, nil
}

func (j *jsonUnmarshaller) Unmarshal(p []byte) (*Entry, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(p, &m); err != nil {
		return nil, err
	}
	id, err := j.getAndDeleteKey(m, j.jsonKeys.id, true)
	if err != nil {
		return nil, err
	}
	timeString, err := j.getAndDeleteKey(m, j.jsonKeys.time, true)
	if err != nil {
		return nil, err
	}
	t, err := time.Parse(timeFormat, timeString.(string))
	if err != nil {
		return nil, err
	}
	levelString, err := j.getAndDeleteKey(m, j.jsonKeys.level, true)
	if err != nil {
		return nil, err
	}
	levelObj, ok := Level_value[strings.ToUpper(levelString.(string))]
	if !ok {
		return nil, fmt.Errorf("ledge: no level for name %s", strings.ToUpper(levelString.(string)))
	}
	level := Level(levelObj)
	writerOutputObj, err := j.getAndDeleteKey(m, j.jsonKeys.writerOutput, false)
	if err != nil {
		return nil, err
	}
	var writerOutput []byte
	switch writerOutputObj.(type) {
	case string:
		writerOutput = []byte(writerOutputObj.(string))
	case []byte:
		writerOutput = writerOutputObj.([]byte)
	default:
		return nil, fmt.Errorf("writerOutput is of type %T", writerOutputObj)
	}
	if len(writerOutput) == 0 {
		writerOutput = nil
	}
	eventType, err := j.getAndDeleteKey(m, j.jsonKeys.eventType, true)
	if err != nil {
		return nil, err
	}
	eventObj, err := j.getAndDeleteKey(m, eventType.(string), true)
	if err != nil {
		return nil, err
	}
	event, err := j.reflectTypeHandler.getEvent(eventType.(string), eventObj)
	if err != nil {
		return nil, err
	}
	var contexts []Context
	for contextType, contextObj := range m {
		context, err := j.reflectTypeHandler.getContext(contextType, contextObj)
		if err != nil {
			return nil, err
		}
		contexts = append(contexts, context)
	}
	return &Entry{
		ID:           id.(string),
		Time:         t,
		Level:        level,
		Contexts:     contexts,
		Event:        Event(event),
		WriterOutput: writerOutput,
	}, nil
}

func (j *jsonUnmarshaller) getAndDeleteKey(m map[string]interface{}, key string, required bool) (interface{}, error) {
	value, ok := m[key]
	if !ok {
		if required {
			return "", fmt.Errorf("ledge: no value for %s", key)
		}
		return "", nil
	}
	delete(m, key)
	return value, nil
}

type protoUnmarshaller struct {
	reflectTypeHandler *reflectTypeHandler
}

func newProtoUnmarshaller(
	specification *Specification,
) (*protoUnmarshaller, error) {
	reflectTypeHandler, err := newReflectTypeHandler(specification)
	if err != nil {
		return nil, err
	}
	return &protoUnmarshaller{
		reflectTypeHandler,
	}, nil
}

func (p *protoUnmarshaller) Unmarshal(buffer []byte) (*Entry, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(buffer))
	bBuffer := bytes.NewBuffer(nil)
	if _, err := bBuffer.ReadFrom(decoder); err != nil {
		return nil, err
	}
	b := bBuffer.Bytes()
	protoEntry := &ProtoEntry{}
	if err := proto.Unmarshal(b, protoEntry); err != nil {
		return nil, err
	}
	entry := &Entry{
		ID:           protoEntry.Id,
		Time:         time.Unix(protoEntry.TimeUnixNsec/int64(time.Second), protoEntry.TimeUnixNsec%int64(time.Second)).UTC(),
		Level:        protoEntry.Level,
		Contexts:     make([]Context, 0),
		WriterOutput: protoEntry.WriterOutput,
	}
	event, err := p.reflectTypeHandler.getEvent(protoEntry.EventTypeName, protoEntry.Event)
	if err != nil {
		return nil, err
	}
	entry.Event = event
	for contextTypeName, contextBytes := range protoEntry.ContextTypeNameToContext {
		context, err := p.reflectTypeHandler.getContext(contextTypeName, contextBytes)
		if err != nil {
			return nil, err
		}
		entry.Contexts = append(entry.Contexts, context)
	}
	return entry, nil
}
