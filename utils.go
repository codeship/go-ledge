package ledge

import (
	"bytes"
	"fmt"
	"reflect"
)

func mergeSpecifications(specifications []*Specification) *Specification {
	contextReflectTypeToContext := make(map[reflect.Type]Context)
	eventReflectTypeToEvent := make(map[reflect.Type]Event)
	for _, specification := range specifications {
		for _, contextType := range specification.ContextTypes {
			contextReflectTypeToContext[reflect.TypeOf(contextType)] = contextType
		}
		for _, eventType := range specification.EventTypes {
			eventReflectTypeToEvent[reflect.TypeOf(eventType)] = eventType
		}
	}
	var contexts []Context
	var events []Event
	for _, context := range contextReflectTypeToContext {
		contexts = append(contexts, context)
	}
	for _, event := range eventReflectTypeToEvent {
		events = append(events, event)
	}
	return &Specification{
		ContextTypes: contexts,
		EventTypes:   events,
	}
}

func getReflectTypeName(reflectType reflect.Type) (string, error) {
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

func checkEntriesEqual(
	entries []*Entry,
	expected []*Entry,
	checkID bool,
	checkTime bool,
) error {
	// TODO(pedge): just change the tests
	for _, elem := range expected {
		elem.Time = elem.Time.UTC()
	}
	if len(expected) != len(entries) {
		expectedStrings := make([]string, len(expected))
		for i, elem := range expected {
			expectedStrings[i] = fmt.Sprintf("%+v", elem)
		}
		entryStrings := make([]string, len(entries))
		for i, elem := range entries {
			entryStrings[i] = fmt.Sprintf("%+v", elem)
		}
		return fmt.Errorf("ledge: expected %v, got %v (length mismatch of %d and %d)", expectedStrings, entryStrings, len(expected), len(entries))
	}
	for i, elem := range expected {
		if err := checkSingleEntriesEqual(elem, entries[i], checkID, checkTime); err != nil {
			return err
		}
	}
	return nil
}

// reflect.DeepEqual does not work on linux for time.Time
func checkSingleEntriesEqual(expected *Entry, actual *Entry, checkID bool, checkTime bool) error {
	if checkID && expected.ID != actual.ID {
		return fmt.Errorf("ledge: expected %+v, got %+v (ID)", expected, actual)
	}
	if checkTime && !expected.Time.Equal(actual.Time) {
		return fmt.Errorf("ledge: expected %+v, got %+v (Time)", expected, actual)
	}
	if expected.Level != actual.Level {
		return fmt.Errorf("ledge: expected %+v, got %+v (Level)", expected, actual)
	}
	if !reflect.DeepEqual(expected.Contexts, actual.Contexts) {
		if expected.Contexts != nil {
			if len(expected.Contexts) != 0 || len(actual.Contexts) != 0 {
				return fmt.Errorf("ledge: expected %+v, got %+v (Contexts)", expected, actual)
			}
			if actual.Contexts == nil {
				return fmt.Errorf("ledge: expected %+v, got %+v (Actual contexts nil, expected not)", expected, actual)
			}
			return fmt.Errorf("ledge: expected %+v, got %+v (Contexts)", expected, actual)
		}
	}
	if !reflect.DeepEqual(expected.Event, actual.Event) {
		return fmt.Errorf("ledge: expected %+v, got %+v (Event)", expected, actual)
	}
	if !reflect.DeepEqual(expected.WriterOutput, actual.WriterOutput) {
		return fmt.Errorf("ledge: expected %+v, got %+v (WriterOutput)", expected, actual)
	}
	return nil
}
