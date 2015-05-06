package ledge

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type jsonUnmarshaller struct {
	reflectTypeProvider *reflectTypeProvider
	jsonKeys            *jsonKeys
}

func newJSONUnmarshaller(
	specification *Specification,
	jsonKeys *jsonKeys,
) (*jsonUnmarshaller, error) {
	reflectTypeProvider, err := newReflectTypeProvider(specification)
	if err != nil {
		return nil, err
	}
	return &jsonUnmarshaller{
		reflectTypeProvider,
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
	level, err := LevelOf(levelString.(string))
	if err != nil {
		return nil, err
	}
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
	event, err := j.getEvent(eventType.(string), eventObj)
	if err != nil {
		return nil, err
	}
	var contexts []Context
	for contextType, contextObj := range m {
		context, err := j.getContext(contextType, contextObj)
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

func (j *jsonUnmarshaller) getContext(objectType string, object interface{}) (interface{}, error) {
	reflectType, err := j.reflectTypeProvider.getContextReflectType(objectType)
	if err != nil {
		return nil, err
	}
	return j.getObject(reflectType, object)
}

func (j *jsonUnmarshaller) getEvent(objectType string, object interface{}) (interface{}, error) {
	reflectType, err := j.reflectTypeProvider.getEventReflectType(objectType)
	if err != nil {
		return nil, err
	}
	return j.getObject(reflectType, object)
}

func (j *jsonUnmarshaller) getObject(reflectType reflect.Type, object interface{}) (interface{}, error) {
	// the below logic is to type the object correctly, do not mess with it
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	objectPtr := reflect.New(reflectType).Interface()
	if err := json.Unmarshal(data, objectPtr); err != nil {
		return nil, err
	}
	return reflect.ValueOf(objectPtr).Elem().Interface(), nil
}
