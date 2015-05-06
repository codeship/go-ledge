package ledge

import (
	"bytes"
	"fmt"
	"reflect"
)

func getFullyQualifiedName(reflectType reflect.Type) (string, error) {
	buffer := bytes.NewBuffer(nil)
	for reflectType.Kind() == reflect.Ptr {
		if _, err := buffer.WriteString("*"); err != nil {
			return "", err
		}
		reflectType = reflectType.Elem()
	}
	pkgPath := reflectType.PkgPath()
	if pkgPath == "" {
		return "", fmt.Errorf("ledge: no package for type %v", reflectType)
	}
	if _, err := buffer.WriteString("\""); err != nil {
		return "", err
	}
	if _, err := buffer.WriteString(pkgPath); err != nil {
		return "", err
	}
	if _, err := buffer.WriteString("\""); err != nil {
		return "", err
	}
	name := reflectType.Name()
	if name == "" {
		return "", fmt.Errorf("ledge: no name for type %v", reflectType)
	}
	if _, err := buffer.WriteString("."); err != nil {
		return "", err
	}
	if _, err := buffer.WriteString(name); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func includeEntry(filters []Filter, entry *Entry) bool {
	if filters == nil || len(filters) == 0 {
		return true
	}
	for _, filter := range filters {
		if !filter.Include((entry)) {
			return false
		}
	}
	return true
}

// reflect.DeepEqual does not work on linux for time.Time
func entriesEqual(one *Entry, two *Entry, checkID bool, checkTime bool) bool {
	if checkID && one.ID != two.ID {
		return false
	}
	if checkTime && !one.Time.Equal(two.Time) {
		return false
	}
	if one.Level != two.Level {
		return false
	}
	if !reflect.DeepEqual(one.Contexts, two.Contexts) {
		return false
	}
	if !reflect.DeepEqual(one.Event, two.Event) {
		return false
	}
	if !reflect.DeepEqual(one.WriterOutput, two.WriterOutput) {
		return false
	}
	return true

}
