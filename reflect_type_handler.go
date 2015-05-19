package ledge

import (
	"encoding/json"
	"reflect"
)

type reflectTypeHandler struct {
	reflectTypeProvider *reflectTypeProvider
}

func newReflectTypeHandler(
	specification *Specification,
) (*reflectTypeHandler, error) {
	reflectTypeProvider, err := newReflectTypeProvider(specification)
	if err != nil {
		return nil, err
	}
	return &reflectTypeHandler{
		reflectTypeProvider,
	}, nil
}

func (r *reflectTypeHandler) getContext(objectType string, object interface{}) (interface{}, error) {
	reflectType, err := r.reflectTypeProvider.getContextReflectType(objectType)
	if err != nil {
		return nil, err
	}
	return r.getObject(reflectType, object)
}

func (r *reflectTypeHandler) getEvent(objectType string, object interface{}) (interface{}, error) {
	reflectType, err := r.reflectTypeProvider.getEventReflectType(objectType)
	if err != nil {
		return nil, err
	}
	return r.getObject(reflectType, object)
}

func (r *reflectTypeHandler) getObject(reflectType reflect.Type, object interface{}) (interface{}, error) {
	if data, ok := object.([]byte); ok {
		objectPtr := reflect.New(reflectType).Interface()
		if err := unmarshalBinary(data, objectPtr); err != nil {
			return nil, err
		}
		return reflect.ValueOf(objectPtr).Elem().Interface(), nil
	}
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
