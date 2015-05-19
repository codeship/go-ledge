package ledge

import (
	"bytes"
	"encoding/base64"
	"time"

	"github.com/golang/protobuf/proto"
)

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
