package ledge

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
)

type protoUnmarshaller struct {
	reflectTypeProvider *reflectTypeProvider
}

func newProtoUnmarshaller(
	specification *Specification,
) (*protoUnmarshaller, error) {
	reflectTypeProvider, err := newReflectTypeProvider(specification)
	if err != nil {
		return nil, err
	}
	return &protoUnmarshaller{
		reflectTypeProvider,
	}, nil
}

func (p *protoUnmarshaller) Unmarshal(buffer []byte) (*Entry, error) {
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewBuffer(buffer))
	bBuffer := bytes.NewBuffer(nil)
	if _, err := bBuffer.ReadFrom(decoder); err != nil {
		return nil, fmt.Errorf("Failed to decode buffer: %s - %s", err.Error(), string(buffer))
	}
	b := bBuffer.Bytes()
	protoEntry := &ProtoEntry{}
	if err := proto.Unmarshal(b, protoEntry); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal protobuf: %s - %s", err.Error(), string(b))
	}
	entry := &Entry{
		ID:           protoEntry.Id,
		Time:         time.Unix(protoEntry.TimeUnixNsec/int64(time.Second), protoEntry.TimeUnixNsec%int64(time.Second)).UTC(),
		Level:        protoEntry.Level,
		Contexts:     make([]Context, 0),
		WriterOutput: protoEntry.WriterOutput,
	}
	event, err := p.getEvent(protoEntry.EventTypeName, protoEntry.Event)
	if err != nil {
		return nil, err
	}
	entry.Event = event
	for contextTypeName, contextBytes := range protoEntry.ContextTypeNameToContext {
		context, err := p.getContext(contextTypeName, contextBytes)
		if err != nil {
			return nil, err
		}
		entry.Contexts = append(entry.Contexts, context)
	}
	return entry, nil
}

func (p *protoUnmarshaller) getContext(objectType string, object []byte) (interface{}, error) {
	reflectType, err := p.reflectTypeProvider.getContextReflectType(objectType)
	if err != nil {
		return nil, err
	}
	return p.getObject(reflectType, object)
}

func (p *protoUnmarshaller) getEvent(objectType string, object []byte) (interface{}, error) {
	reflectType, err := p.reflectTypeProvider.getEventReflectType(objectType)
	if err != nil {
		return nil, err
	}
	return p.getObject(reflectType, object)
}

func (p *protoUnmarshaller) getObject(reflectType reflect.Type, object []byte) (interface{}, error) {
	if reflectType.Implements(reflect.TypeOf((*proto.Message)(nil)).Elem()) {
		protoMessage := reflect.New(reflectType.Elem()).Interface().(proto.Message)
		if err := proto.Unmarshal(object, protoMessage); err != nil {
			return nil, err
		}
		return protoMessage, nil
	}
	objectPtr := reflect.New(reflectType).Interface()
	if err := gob.NewDecoder(bytes.NewBuffer(object)).Decode(objectPtr); err != nil {
		return nil, err
	}
	return reflect.ValueOf(objectPtr).Elem().Interface(), nil
}
