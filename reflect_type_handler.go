package ledge

import "reflect"

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

func (r *reflectTypeHandler) getContext(objectType string, object []byte) (interface{}, error) {
	reflectType, err := r.reflectTypeProvider.getContextReflectType(objectType)
	if err != nil {
		return nil, err
	}
	return r.getObject(reflectType, object)
}

func (r *reflectTypeHandler) getEvent(objectType string, object []byte) (interface{}, error) {
	reflectType, err := r.reflectTypeProvider.getEventReflectType(objectType)
	if err != nil {
		return nil, err
	}
	return r.getObject(reflectType, object)
}

func (r *reflectTypeHandler) getObject(reflectType reflect.Type, object []byte) (interface{}, error) {
	objectPtr := reflect.New(reflectType).Interface()
	if err := unmarshalBinary(object, objectPtr); err != nil {
		return nil, err
	}
	return reflect.ValueOf(objectPtr).Elem().Interface(), nil
}
