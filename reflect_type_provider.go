package ledge

import (
	"fmt"
	"reflect"
)

type reflectTypeProvider struct {
	contextKeyToReflectType map[string]reflect.Type
	eventKeyToReflectType   map[string]reflect.Type
	contextReflectTypes     map[reflect.Type]bool
	eventReflectTypes       map[reflect.Type]bool
}

func newReflectTypeProvider(
	specification *Specification,
) (*reflectTypeProvider, error) {
	contextKeyToReflectType := make(map[string]reflect.Type)
	eventKeyToReflectType := make(map[string]reflect.Type)
	contextReflectTypes := make(map[reflect.Type]bool)
	eventReflectTypes := make(map[reflect.Type]bool)
	if specification != nil {
		if err := addToKeyToReflectType(contextKeyToReflectType, contextReflectTypes, specification.ContextTypes); err != nil {
			return nil, err
		}
		if err := addToKeyToReflectType(eventKeyToReflectType, eventReflectTypes, specification.EventTypes); err != nil {
			return nil, err
		}
	}
	if err := addToKeyToReflectType(eventKeyToReflectType, eventReflectTypes, DefaultEventTypes); err != nil {
		return nil, err
	}
	return &reflectTypeProvider{
		contextKeyToReflectType,
		eventKeyToReflectType,
		contextReflectTypes,
		eventReflectTypes,
	}, nil
}

func (r *reflectTypeProvider) getContextReflectType(key string) (reflect.Type, error) {
	return r.getReflectType(r.contextKeyToReflectType, key)
}

func (r *reflectTypeProvider) getEventReflectType(key string) (reflect.Type, error) {
	return r.getReflectType(r.eventKeyToReflectType, key)
}

func (r *reflectTypeProvider) validateContextReflectType(reflectType reflect.Type) error {
	return r.validateReflectType(r.contextReflectTypes, reflectType)
}

func (r *reflectTypeProvider) validateEventReflectType(reflectType reflect.Type) error {
	return r.validateReflectType(r.eventReflectTypes, reflectType)
}

func (r *reflectTypeProvider) getReflectType(m map[string]reflect.Type, key string) (reflect.Type, error) {
	reflectType, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("ledge: no reflect type for %s", key)
	}
	return reflectType, nil
}

func (r *reflectTypeProvider) validateReflectType(m map[reflect.Type]bool, reflectType reflect.Type) error {
	if _, ok := m[reflectType]; !ok {
		return fmt.Errorf("ledge: reflect type %s not part of specification", reflectType)
	}
	return nil
}

func addToKeyToReflectType(keyToReflectType map[string]reflect.Type, reflectTypes map[reflect.Type]bool, types []interface{}) error {
	for _, t := range types {
		reflectType := reflect.TypeOf(t)
		key, err := getFullyQualifiedName(reflectType)
		if err != nil {
			return err
		}
		keyToReflectType[key] = reflectType
		reflectTypes[reflectType] = true
	}
	return nil
}
